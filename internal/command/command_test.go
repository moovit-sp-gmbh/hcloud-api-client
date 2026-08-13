package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteReturnsResult(t *testing.T) {
	cmd := New(&cobra.Command{
		Use: "greet",
		RunE: func(cmd *cobra.Command, args []string) error {
			SetResult(cmd, Text("hello"))
			return nil
		},
	})

	out, err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected a result, got nil")
	}
	if got := out.String(); got != "hello" {
		t.Fatalf("result = %q, want %q", got, "hello")
	}
}

func TestExecuteCapturesInCommandHelp(t *testing.T) {
	cmd := New(&cobra.Command{
		Use:   "group",
		Short: "a command group",
		Long:  "group is a parent command",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	})

	out, err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected captured help output, got nil")
	}
	if !strings.Contains(out.String(), "group is a parent command") {
		t.Fatalf("captured output missing help text:\n%s", out.String())
	}
}

func TestExecutePropagatesError(t *testing.T) {
	wantErr := errors.New("boom")
	cmd := New(&cobra.Command{
		Use:           "fail",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return wantErr
		},
	})

	out, err := cmd.Execute(nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
	if out != nil {
		t.Fatalf("expected nil result on error, got %v", out)
	}
}
