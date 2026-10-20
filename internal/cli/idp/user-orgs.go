package idp

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/idpClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	userOrgsSearchCommand.Flags().String("token", "", "authorization token")
	userOrgsSearchCommand.Flags().Float32("limit", 0, "Page size")
	userOrgsSearchCommand.Flags().Float32("page", 0, "Page index")
	userOrgsSearchCommand.Flags().Bool("teamsOfUser", false, "If true, an additional property holding all teams that the requesting user is part of and that belong to the organization will be returned")
	userOrgsSearchCommand.Flags().Bool("totalMemberCount", false, "If true, an additional property holding the total member count of the organization will be returned")
	userOrgsSearchCommand.Flags().Float32("membersSample", 0, "If set, an additional property holding a sample of members will be returned. The provided number represents the max size of the sample")
	userOrgsSearchCommand.Flags().String("body", "{}", "request body as JSON")
}

var userOrgsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search for organizations with search filters",
	Long:  "Required scopes: IDP_ORGANIZATION_READ",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewIDPClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &idpClient.PostV1UserOrgsSearchParams{}

		if cmd.Flags().Changed("limit") {
			limit, err := cmd.Flags().GetFloat32("limit")
			if err != nil {
				return err
			}
			params.Limit = &limit
		}

		if cmd.Flags().Changed("page") {
			page, err := cmd.Flags().GetFloat32("page")
			if err != nil {
				return err
			}
			params.Page = &page
		}

		if cmd.Flags().Changed("teamsOfUser") {
			teamsOfUser, err := cmd.Flags().GetBool("teamsOfUser")
			if err != nil {
				return err
			}
			params.TeamsOfUser = &teamsOfUser
		}

		if cmd.Flags().Changed("totalMemberCount") {
			totalMemberCount, err := cmd.Flags().GetBool("totalMemberCount")
			if err != nil {
				return err
			}
			params.TotalMemberCount = &totalMemberCount
		}

		if cmd.Flags().Changed("membersSample") {
			membersSample, err := cmd.Flags().GetFloat32("membersSample")
			if err != nil {
				return err
			}
			params.MembersSample = &membersSample
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1UserOrgsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1UserOrgsSearchWithResponse(context.TODO(), params, body, func(_ context.Context, req *http.Request) error {
			req.Header.Set("authorization", jwt)
			return nil
		})
		if err != nil {
			return err
		}

		var data any
		err = json.Unmarshal(res.Body, &data)
		if err != nil {
			return err
		}
		text, err := command.JSON(data)
		if err != nil {
			return err
		}

		if res.StatusCode() >= 200 && res.StatusCode() < 300 {
			service.Push(latest)
		}

		command.SetResult(cmd, text)

		return nil
	},
}
