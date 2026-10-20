package cosmo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"
	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/command"
	"github.com/spf13/cobra"
)

func init() {
	spacesNamespacesCommentsDeleteCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsDeleteCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsDeleteCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsDeleteCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsDeleteCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsDeleteCommand.Flags().String("body", "{}", "request body as JSON")
	spacesNamespacesCommentsListCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsListCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsListCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsListCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsListCommand.Flags().String("refId", "", "ID of the entity the comments should refer to")
	spacesNamespacesCommentsListCommand.Flags().Float32("page", 0, "Page index")
	spacesNamespacesCommentsListCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsListCommand.Flags().Float32("replySample", 0, "Page size")
	spacesNamespacesCommentsListCommand.Flags().String("sort_field", "", "Field to sort comments by. NEWEST/OLDEST sort by create date desc/asc.")
	spacesNamespacesCommentsListCommand.Flags().String("sort_direction", "", "Sort direction for comments (ignored for NEWEST/OLDEST).")
	spacesNamespacesCommentsCreateCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsCreateCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsCreateCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsCreateCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsCreateCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsCreateCommand.Flags().String("refId", "", "ID of the item referenced")
	spacesNamespacesCommentsCreateCommand.Flags().String("text", "", "Comment text")
	spacesNamespacesCommentsCreateCommand.Flags().String("bodyAnnotation", "", "Optional annotation for the comment")
	spacesNamespacesCommentsCreateCommand.Flags().String("color", "", "Comment color")
	spacesNamespacesCommentsCreateCommand.Flags().StringSlice("mentions", nil, "Array of mentioned user IDs")
	spacesNamespacesCommentsCreateCommand.Flags().StringSlice("mentionsTeams", nil, "Array of mentioned team ids")
	spacesNamespacesCommentsCreateCommand.Flags().String("visibility", "", "Visibility of the comment. Internal comments are only visible to org members.")
	spacesNamespacesCommentsSearchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsSearchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsSearchCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsSearchCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsSearchCommand.Flags().String("refId", "", "ID of the entity the comments should refer to")
	spacesNamespacesCommentsSearchCommand.Flags().Float32("limit", 0, "Page size")
	spacesNamespacesCommentsSearchCommand.Flags().Float32("page", 0, "Page index")
	spacesNamespacesCommentsSearchCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsSearchCommand.Flags().String("body", "{}", "request body as JSON")
	spacesNamespacesCommentsGetCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsGetCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsGetCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsGetCommand.Flags().String("commentId", "", "")
	spacesNamespacesCommentsGetCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsGetCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsPatchCommand.Flags().String("orgName", "", "Name of the organization (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsPatchCommand.Flags().String("spaceName", "", "Name of the space (only alphanumeric characters, hyphens and dots are allowed; it should contain between 1 and 32 characters)")
	spacesNamespacesCommentsPatchCommand.Flags().String("namespaceName", "", "Name of the namespace (only alphanumeric characters, hyphens and dots are allowed; it should contain between 3 and 64 characters)")
	spacesNamespacesCommentsPatchCommand.Flags().String("commentId", "", "")
	spacesNamespacesCommentsPatchCommand.Flags().String("token", "", "authorization token")
	spacesNamespacesCommentsPatchCommand.Flags().Bool("annotation", false, "")
	spacesNamespacesCommentsPatchCommand.Flags().String("bodyAnnotation", "", "Optional annotation for the comment")
	spacesNamespacesCommentsPatchCommand.Flags().String("color", "", "Comment color")
	spacesNamespacesCommentsPatchCommand.Flags().Bool("completed", false, "Set the comment to completed/uncompleted.")
	spacesNamespacesCommentsPatchCommand.Flags().StringSlice("mentions", nil, "Array of mentioned user IDs")
	spacesNamespacesCommentsPatchCommand.Flags().StringSlice("mentionsTeams", nil, "Array of mentioned team ids")
	spacesNamespacesCommentsPatchCommand.Flags().String("text", "", "Comment text")
	spacesNamespacesCommentsPatchCommand.Flags().String("visibility", "", "Visibility of the comment. Internal comments are only visible to org members.")

	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsDeleteCommand)
	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsListCommand)
	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsCreateCommand)
	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsSearchCommand)
	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsGetCommand)
	spacesNamespacesCommentsCommand.AddCommand(spacesNamespacesCommentsPatchCommand)
}

var spacesNamespacesCommentsCommand = &cobra.Command{
	Use:   "comments",
	Short: "Commands for comments",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var spacesNamespacesCommentsDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete Comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsParams{}

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.DeleteV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.DeleteV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsWithResponse(context.TODO(), orgName, spaceName, namespaceName, params, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesCommentsListCommand = &cobra.Command{
	Use:   "list",
	Short: "List comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		refId, err := command.FlagOrPromptString(cmd, "refId")
		if err != nil {
			return err
		}
		refIdParsed, err := uuid.Parse(refId)
		if err != nil {
			return err
		}

		params := &cosmoClient.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsParams{
			RefId: refIdParsed,
		}

		if cmd.Flags().Changed("page") {
			page, err := cmd.Flags().GetFloat32("page")
			if err != nil {
				return err
			}
			params.Page = &page
		}

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		if cmd.Flags().Changed("replySample") {
			replySample, err := cmd.Flags().GetFloat32("replySample")
			if err != nil {
				return err
			}
			params.ReplySample = &replySample
		}

		if cmd.Flags().Changed("sort_field") {
			sortField, err := cmd.Flags().GetString("sort_field")
			if err != nil {
				return err
			}
			sortFieldValue := cosmoClient.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsParamsSortField(sortField)
			params.SortField = &sortFieldValue
		}

		if cmd.Flags().Changed("sort_direction") {
			sortDirection, err := cmd.Flags().GetString("sort_direction")
			if err != nil {
				return err
			}
			sortDirectionValue := cosmoClient.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsParamsSortDirection(sortDirection)
			params.SortDirection = &sortDirectionValue
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsWithResponse(context.TODO(), orgName, spaceName, namespaceName, params, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesCommentsCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create Comment",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsParams{}

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		refId, err := command.FlagOrPromptString(cmd, "refId")
		if err != nil {
			return err
		}
		refIdParsed, err := uuid.Parse(refId)
		if err != nil {
			return err
		}

		textValue, err := command.FlagOrPromptString(cmd, "text")
		if err != nil {
			return err
		}

		body := cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsJSONRequestBody{
			RefId: refIdParsed,
			Text:  textValue,
		}

		if cmd.Flags().Changed("bodyAnnotation") {
			bodyAnnotation, err := cmd.Flags().GetString("bodyAnnotation")
			if err != nil {
				return err
			}
			body.Annotation = &bodyAnnotation
		}

		if cmd.Flags().Changed("color") {
			color, err := cmd.Flags().GetString("color")
			if err != nil {
				return err
			}
			colorValue := cosmoClient.CommentCreateDtoColor(color)
			body.Color = &colorValue
		}

		if cmd.Flags().Changed("mentions") {
			mentions, err := cmd.Flags().GetStringSlice("mentions")
			if err != nil {
				return err
			}
			body.Mentions = &mentions
		}

		if cmd.Flags().Changed("mentionsTeams") {
			mentionsTeams, err := cmd.Flags().GetStringSlice("mentionsTeams")
			if err != nil {
				return err
			}
			body.MentionsTeams = &mentionsTeams
		}

		if cmd.Flags().Changed("visibility") {
			visibility, err := cmd.Flags().GetString("visibility")
			if err != nil {
				return err
			}
			visibilityValue := cosmoClient.CommentCreateDtoVisibility(visibility)
			body.Visibility = &visibilityValue
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsWithResponse(context.TODO(), orgName, spaceName, namespaceName, params, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesCommentsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		refId, err := command.FlagOrPromptString(cmd, "refId")
		if err != nil {
			return err
		}
		refIdParsed, err := uuid.Parse(refId)
		if err != nil {
			return err
		}

		params := &cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsSearchParams{
			RefId: refIdParsed,
		}

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

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		bodyJSON, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var body cosmoClient.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsSearchJSONRequestBody
		err = json.Unmarshal([]byte(bodyJSON), &body)
		if err != nil {
			return err
		}

		res, err := client.PostV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsSearchWithResponse(context.TODO(), orgName, spaceName, namespaceName, params, body, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesCommentsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get single comment by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		commentId, err := command.FlagOrPromptString(cmd, "commentId")
		if err != nil {
			return err
		}
		commentIdParsed, err := uuid.Parse(commentId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsCommentIdParams{}

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		res, err := client.GetV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsCommentIdWithResponse(context.TODO(), orgName, spaceName, namespaceName, commentIdParsed, params, func(_ context.Context, req *http.Request) error {
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

var spacesNamespacesCommentsPatchCommand = &cobra.Command{
	Use:   "patch",
	Short: "Edit Comment",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clients.NewCosmoClient(cmd)
		if err != nil {
			return err
		}

		service := internal.New()
		latest, err := service.GetLatest()
		if err != nil {
			return err
		}

		orgName, err := command.GetArg(cmd, command.StringFlagIfSet("orgName"), command.ContextIfValidString(latest.OrgName), command.PromptString("orgName"))
		if err != nil {
			return err
		}
		latest.OrgName = sql.NullString{Valid: true, String: orgName}

		spaceName, err := command.FlagOrPromptString(cmd, "spaceName")
		if err != nil {
			return err
		}

		namespaceName, err := command.FlagOrPromptString(cmd, "namespaceName")
		if err != nil {
			return err
		}

		commentId, err := command.FlagOrPromptString(cmd, "commentId")
		if err != nil {
			return err
		}
		commentIdParsed, err := uuid.Parse(commentId)
		if err != nil {
			return err
		}

		jwt, err := command.GetArg(cmd, command.StringFlagIfSet("token"), command.ContextIfValidString(latest.Jwt), command.PromptString("authorization token"))
		if err != nil {
			return err
		}
		latest.Jwt = sql.NullString{Valid: true, String: jwt}

		params := &cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsCommentIdParams{}

		if cmd.Flags().Changed("annotation") {
			annotation, err := cmd.Flags().GetBool("annotation")
			if err != nil {
				return err
			}
			params.Annotation = &annotation
		}

		body := cosmoClient.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsCommentIdJSONRequestBody{}

		if cmd.Flags().Changed("bodyAnnotation") {
			bodyAnnotation, err := cmd.Flags().GetString("bodyAnnotation")
			if err != nil {
				return err
			}
			body.Annotation = &bodyAnnotation
		}

		if cmd.Flags().Changed("color") {
			color, err := cmd.Flags().GetString("color")
			if err != nil {
				return err
			}
			colorValue := cosmoClient.CommentEditDtoColor(color)
			body.Color = &colorValue
		}

		if cmd.Flags().Changed("completed") {
			completed, err := cmd.Flags().GetBool("completed")
			if err != nil {
				return err
			}
			body.Completed = &completed
		}

		if cmd.Flags().Changed("mentions") {
			mentions, err := cmd.Flags().GetStringSlice("mentions")
			if err != nil {
				return err
			}
			body.Mentions = &mentions
		}

		if cmd.Flags().Changed("mentionsTeams") {
			mentionsTeams, err := cmd.Flags().GetStringSlice("mentionsTeams")
			if err != nil {
				return err
			}
			body.MentionsTeams = &mentionsTeams
		}

		if cmd.Flags().Changed("text") {
			textValue, err := cmd.Flags().GetString("text")
			if err != nil {
				return err
			}
			body.Text = &textValue
		}

		if cmd.Flags().Changed("visibility") {
			visibility, err := cmd.Flags().GetString("visibility")
			if err != nil {
				return err
			}
			visibilityValue := cosmoClient.CommentEditDtoVisibility(visibility)
			body.Visibility = &visibilityValue
		}

		res, err := client.PatchV1OrgOrgNameSpacesSpaceNameNamespacesNamespaceNameCommentsCommentIdWithResponse(context.TODO(), orgName, spaceName, namespaceName, commentIdParsed, params, body, func(_ context.Context, req *http.Request) error {
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
