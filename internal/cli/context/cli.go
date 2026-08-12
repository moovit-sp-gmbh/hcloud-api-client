package context

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/db"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().Int64P("limit", "l", 5, "number of results per page")
	listCmd.Flags().Int64P("page", "p", 0, "page offset")
}

var listCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		service := internal.New()

		page, err := cmd.Flags().GetInt64("page")
		if err != nil {
			return err
		}

		limit, err := cmd.Flags().GetInt64("limit")
		if err != nil {
			return err
		}

		ctxs, err := service.List(page, limit)
		if err != nil {
			return err
		} else if ctxs == nil {
			ctxs = []db.Context{}
		}

		data, err := json.MarshalIndent(ctxs, "", "  ")
		if err != nil {
			return err
		}

		command.SetResult(cmd, command.Text(string(data)))
		return nil
	},
}

var latestCmd = &cobra.Command{
	Use: "latest",
	RunE: func(cmd *cobra.Command, args []string) error {
		service := internal.New()

		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		text, err := command.JSON(ctx)
		if err != nil {
			return err
		}

		command.SetResult(cmd, text)
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set a context as the latest",
	Example: "context set <context ID>",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid context ID: %w", err)
		}

		service := internal.New()

		err = service.SetLatest(id)
		if err != nil {
			return err
		}

		command.SetResult(cmd, command.Text(fmt.Sprintf("Set context %d as latest", id)))
		return nil
	},
}

func Command() *cobra.Command {
	rootCmd := &cobra.Command{Use: "context", Run: func(cmd *cobra.Command, args []string) { cmd.Help() }}
	rootCmd.Groups()
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(latestCmd)
	rootCmd.AddCommand(setCmd)
	return rootCmd
}
