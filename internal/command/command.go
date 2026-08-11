package command

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Executor interface {
	Execute([]string) (fmt.Stringer, error)
	Traverse([]string) (Executor, []string, error)
	Help([]string) []string
}

type Text string

func JSON(obj any) (Text, error) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return Text(string(data)), nil
}

func (t Text) String() string { return string(t) }

// resultKey is the private context key under which a command's result holder
// is carried during execution.
type _resultKey struct{}

var resultKey = _resultKey{}

// holder is a mutable box a command writes its result into via [SetResult].
type holder struct {
	value fmt.Stringer
}

type Command struct {
	*cobra.Command
}

func New(cmd *cobra.Command) *Command {
	return &Command{Command: cmd}
}

func (c *Command) Execute(args []string) (fmt.Stringer, error) {
	h := &holder{}
	ctx := context.WithValue(context.Background(), resultKey, h)

	// Point the command's output at a buffer so anything it writes itself —
	// most notably an in-command cmd.Help() call on a command group — is
	// captured. OutOrStdout() walks up to the parent, so setting it on the
	// root also captures writes from the traversed subcommand.
	var buf bytes.Buffer
	c.Command.SetArgs(args)
	c.Command.SetOut(&buf)
	defer c.Command.SetOut(nil)

	if err := c.Command.ExecuteContext(ctx); err != nil {
		return nil, err
	}

	// An explicit SetResult wins; otherwise fall back to captured output.
	if h.value != nil {
		return h.value, nil
	}
	if buf.Len() > 0 {
		return Text(buf.String()), nil
	}
	return nil, nil
}

func SetResult(cmd *cobra.Command, value fmt.Stringer) {
	if h, ok := cmd.Context().Value(resultKey).(*holder); ok {
		h.value = value
	}
}

func (c *Command) Traverse(args []string) (Executor, []string, error) {
	cmd, rest, err := c.Command.Traverse(args)
	if err != nil {
		return nil, rest, err
	}
	return New(cmd), rest, nil
}

func (c *Command) Help(args []string) []string {
	var items []string

	flags := map[string]*pflag.Flag{}
	shorthands := map[string]*pflag.Flag{}
	c.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
		if f.Shorthand != "" {
			shorthands[f.Shorthand] = f
		}
	})
	c.Flags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
		if f.Shorthand != "" {
			shorthands[f.Shorthand] = f
		}
	})

	var lastArg string
	if len(args) > 0 {
		lastArg = args[len(args)-1]
	}
	if len(lastArg) > 1 && lastArg[0] == '-' && lastArg[1] == '-' {
		if f, ok := flags[lastArg[2:]]; ok {
			if f.Value.Type() != "bool" {
				return nil
			}
		}

		partialMatches := findPartials(flags, lastArg)
		if len(partialMatches) > 0 {
			for i := range partialMatches {
				partialMatches[i] = "--" + partialMatches[i]
			}
			return partialMatches
		}
	} else if len(lastArg) > 0 && lastArg[0] == '-' {
		if f, ok := shorthands[lastArg[1:]]; ok {
			if f.Value.Type() != "bool" {
				return nil
			}
		}
		partialMatches := findPartials(shorthands, lastArg)
		if len(partialMatches) > 0 {
			for i := range partialMatches {
				partialMatches[i] = "-" + partialMatches[i]
			}
			return partialMatches
		}
	}

	subcmds := c.Commands()
	for _, subcmd := range subcmds {
		if strings.HasPrefix(subcmd.Use, lastArg) {
			items = append(items, subcmd.Use)
		}
	}

	if len(items) == 0 {
		for _, subcmd := range subcmds {
			items = append(items, subcmd.Use)
		}
		for flag := range flags {
			items = append(items, fmt.Sprintf("--%s", flag))
		}
	}

	return items
}

func findPartials(m map[string]*pflag.Flag, start string) []string {
	var ret []string
	for k := range m {
		if strings.HasPrefix(k, start) {
			ret = append(ret, k)
		}
	}
	return ret
}

func Server(cmd *cobra.Command) (string, error) {
	service := internal.New()
	ctx, err := service.GetLatest()
	if err != nil {
		return "", err
	}

	serverSet := cmd.Flags().Changed("server")
	log.Println(serverSet)
	if !serverSet {
		if ctx.Server.Valid {
			return ctx.Server.String, nil
		} else {
			serverFlag, err := cmd.Flags().GetString("server")
			if err != nil {
				return "", err
			}
			return serverFlag, nil
		}
	} else {
		serverFlag, err := cmd.Flags().GetString("server")
		if err == nil {
			ctx.Server = sql.NullString{Valid: true, String: serverFlag}
			service.Push(ctx)
		}
		return serverFlag, err
	}
}
