package ui

import "github.com/derailed/tview"

type Command struct {
	*tview.TextView
}

// NewCommand returns a new command.
func NewCommand() *Command {
	c := Command{
		TextView: tview.NewTextView(),
	}

	return &c
}
