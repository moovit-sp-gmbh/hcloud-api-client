package context

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
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

		err = json.NewEncoder(os.Stdout).Encode(ctxs)
		if err != nil {
			return err
		}

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

		return nil
	},
}

func Command() *cobra.Command {
	rootCmd := &cobra.Command{Use: "context", Run: func(cmd *cobra.Command, args []string) { cmd.Help() }}
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(setCmd)
	return rootCmd
}
