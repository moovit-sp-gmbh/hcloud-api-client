package idp

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	orgCommand.Flags().StringP("name", "n", "", "Name of the organization")

	createOrgCommand.Flags().String("company", "", "Company name")

	orgCommand.AddCommand(createOrgCommand)
}

var orgCommand = &cobra.Command{
	Use:   "organization",
	Short: "Organization entrypoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var createOrgCommand = &cobra.Command{
	Use:   "create <name>",
	Short: "Create an organization",
	Long:  "Create an organization with a given name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgName := args[0]
		company, err := cmd.Flags().GetString("company")
		if err != nil {
			return err
		}

		service := internal.New()
		ctx, err := service.GetLatest()
		if err != nil {
			return err
		}

		if !ctx.Jwt.Valid {
			return errors.New("JWT not set. Login first")
		}

		body := map[string]string{
			"name": orgName,
		}

		if company != "" {
			body["company"] = company
		}

		url, err := util.BuildURL(server, "/api/account/v1/org")
		if err != nil {
			return err
		}
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", ctx.Jwt.String)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode >= 300 {
			text, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				return err
			}
			return errors.New(string(text))
		}

		text, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}
		fmt.Println(string(text))

		ctx.OrgName = sql.NullString{Valid: true, String: orgName}
		service.Push(ctx)

		return nil
	},
}
