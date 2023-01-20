package ui

import (
	"github.com/derailed/tview"
)

const (
	defaultSpacer = 4
)

// Prompt captures users free from command input.
type Prompt struct {
	*tview.TextView

	app     *App
	noIcons bool
	icon    rune
	spacer  int
}

// NewPrompt returns a new command view.
func NewPrompt(app *App) *Prompt {
	p := Prompt{
		app:      app,
		TextView: tview.NewTextView(),
		spacer:   defaultSpacer,
	}

	p.SetWordWrap(true)
	p.SetWrap(true)
	p.SetDynamicColors(true)
	p.SetBorder(true)
	p.SetBorderPadding(0, 0, 1, 1)

	return &p
}
