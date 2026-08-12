package command

import (
	"strconv"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli"
	"github.com/spf13/cobra"
)

func FlagOrPromptString(cmd *cobra.Command, name string) (string, error) {
	return FlagOrDifferentPromptString(cmd, name, name)
}

func FlagOrPromptBool(cmd *cobra.Command, name string) (bool, error) {
	return FlagOrDifferentPromptBool(cmd, name, name)
}

func FlagOrPromptInt(cmd *cobra.Command, name string) (int, error) {
	return FlagOrDifferentPromptInt(cmd, name, name)
}

func FlagOrPromptFloat(cmd *cobra.Command, name string) (float64, error) {
	return FlagOrDifferentPromptFloat(cmd, name, name)
}

func FlagOrDifferentPromptString(cmd *cobra.Command, name string, prompt string) (string, error) {
	value, err := cmd.Flags().GetString(name)
	if err == nil && cmd.Flags().Changed(name) {
		return value, nil
	}
	return cli.Prompt(prompt)
}

func FlagOrDifferentPromptBool(cmd *cobra.Command, name string, prompt string) (bool, error) {
	value, err := cmd.Flags().GetBool(name)
	if err == nil && cmd.Flags().Changed(name) {
		return value, nil
	}
	if value, err := cli.Prompt(prompt); err != nil {
		return false, err
	} else {
		return value == "true", nil
	}
}

func FlagOrDifferentPromptInt(cmd *cobra.Command, name string, prompt string) (int, error) {
	value, err := cmd.Flags().GetInt(name)
	if err == nil && cmd.Flags().Changed(name) {
		return value, nil
	}
	str, err := cli.Prompt(prompt)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func FlagOrDifferentPromptFloat(cmd *cobra.Command, name string, prompt string) (float64, error) {
	value, err := cmd.Flags().GetFloat64(name)
	if err == nil && cmd.Flags().Changed(name) {
		return value, nil
	}
	str, err := cli.Prompt(prompt)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(str, 64)
}
