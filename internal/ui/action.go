package ui

import "github.com/gdamore/tcell/v2"

type (
	// ActionHandler handles a keyboard command.
	ActionHandler func(*tcell.EventKey) *tcell.EventKey

	// KeyAction represents a keyboard action.
	KeyAction struct {
		Key         string
		Description string
		Action      ActionHandler
		Visible     bool
	}

	// KeyActions tracks mappings between keystrokes and actions.
	KeyActions map[tcell.Key]KeyAction
)

// NewKeyAction returns a new keyboard action.
func NewKeyAction(k string, d string, a ActionHandler, visible bool) KeyAction {
	return KeyAction{Key: k, Description: d, Action: a, Visible: visible}
}
