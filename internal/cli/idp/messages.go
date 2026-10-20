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
	messagesGetCommand.Flags().String("messageId", "", "")
	messagesGetCommand.Flags().String("token", "", "authorization token")
	messagesBulkReadCommand.Flags().String("token", "", "authorization token")
	messagesBulkReadCommand.Flags().StringSlice("ids", nil, "IDs of the messages to update")
	messagesBulkReadCommand.Flags().Bool("read", false, "Read status to apply to all specified messages")
	messagesSearchCommand.Flags().String("token", "", "authorization token")
	messagesSearchCommand.Flags().Float32("limit", 0, "Page size")
	messagesSearchCommand.Flags().Float32("page", 0, "Page index")
	messagesSearchCommand.Flags().String("body", "{}", "request body as JSON")
	messagesDeleteCommand.Flags().String("messageId", "", "")
	messagesDeleteCommand.Flags().String("token", "", "authorization token")
	messagesDeletedCommand.Flags().String("messageId", "", "")
	messagesDeletedCommand.Flags().String("token", "", "authorization token")
	messagesDeletedCommand.Flags().Bool("deleted", false, "Indicates whether the message is soft-deleted")
	messagesReadCommand.Flags().String("messageId", "", "")
	messagesReadCommand.Flags().String("token", "", "authorization token")
	messagesReadCommand.Flags().Bool("read", false, "Indicates whether the message has been read")
	messagesCommand.AddCommand(messagesGetCommand)
	messagesCommand.AddCommand(messagesBulkReadCommand)
	messagesCommand.AddCommand(messagesSearchCommand)
	messagesCommand.AddCommand(messagesDeleteCommand)
	messagesCommand.AddCommand(messagesDeletedCommand)
	messagesCommand.AddCommand(messagesReadCommand)
}

var messagesCommand = &cobra.Command{
	Use:   "messages",
	Short: "Commands for messages",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var messagesGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieves a message by it's ID",
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

		messageId, err := command.FlagOrDifferentPromptString(cmd, "messageId", "message ID")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.GetV1MessagesMessageIdWithResponse(context.TODO(), messageId, func(_ context.Context, req *http.Request) error {
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

var messagesBulkReadCommand = &cobra.Command{
	Use:   "bulk-read",
	Short: "Bulk-update the read status of multiple messages by ID",
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

		ids, err := cmd.Flags().GetStringSlice("ids")
		if err != nil {
			return err
		}

		read, err := command.FlagOrPromptBool(cmd, "read")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1MessagesBulkReadJSONRequestBody{
			Ids:  ids,
			Read: read,
		}

		res, err := client.PatchV1MessagesBulkReadWithResponse(context.TODO(), body, func(_ context.Context, req *http.Request) error {
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

var messagesSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Retrieves all messages to the specified user that match the provided search filter. Will return all user messages if no filter provided",
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

		params := &idpClient.PostV1MessagesSearchParams{}

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

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body idpClient.PostV1MessagesSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1MessagesSearchWithResponse(context.TODO(), params, body, func(_ context.Context, req *http.Request) error {
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

var messagesDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete message by ID",
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

		messageId, err := command.FlagOrPromptString(cmd, "messageId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		res, err := client.DeleteV1MessagesMessageIdWithResponse(context.TODO(), messageId, func(_ context.Context, req *http.Request) error {
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

var messagesDeletedCommand = &cobra.Command{
	Use:   "deleted",
	Short: "Set the soft-delete state of a received message",
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

		messageId, err := command.FlagOrPromptString(cmd, "messageId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		deleted, err := command.FlagOrPromptBool(cmd, "deleted")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1MessagesMessageIdDeletedJSONRequestBody{
			Deleted: deleted,
		}

		res, err := client.PatchV1MessagesMessageIdDeletedWithResponse(context.TODO(), messageId, body, func(_ context.Context, req *http.Request) error {
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

var messagesReadCommand = &cobra.Command{
	Use:   "read",
	Short: "Set the status of the received message as read or unread",
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

		messageId, err := command.FlagOrPromptString(cmd, "messageId")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		read, err := command.FlagOrPromptBool(cmd, "read")
		if err != nil {
			return err
		}

		body := idpClient.PatchV1MessagesMessageIdReadJSONRequestBody{
			Read: read,
		}

		res, err := client.PatchV1MessagesMessageIdReadWithResponse(context.TODO(), messageId, body, func(_ context.Context, req *http.Request) error {
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
