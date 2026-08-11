package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli/context"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli/idp"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/repl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "Run the command and stay in the REPL")
	rootCmd.AddCommand(idp.Command())
	rootCmd.AddCommand(context.Command())
}

var rootCmd = &cobra.Command{
	Use:   "hcloud",
	Short: "Command-Line Interface for helmut.cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive, err := cmd.PersistentFlags().GetBool("interactive")
		if err != nil {
			return err
		}
		if interactive {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}
			defer f.Close()

			cmd.PersistentFlags().Set("interactive", "false")
			p := tea.NewProgram(repl.New(command.New(cmd), nil))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		} else {
			cmd.Help()
		}
		return nil
	},
}

func Execute() {
	out, err := command.New(rootCmd).Execute(os.Args[1:])
	if err != nil {
		// cobra already prints the error to stderr by default.
		os.Exit(1)
	}
	if out != nil {
		if s := out.String(); s != "" {
			if strings.HasSuffix(s, "\n") {
				fmt.Fprint(os.Stdout, s)
			} else {
				fmt.Fprintln(os.Stdout, s)
			}
		}
	}
}
