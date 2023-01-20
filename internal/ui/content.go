package ui

import (
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
)

var activatedCallback map[string]func()

type Content struct {
	*tview.Pages
}

// NewContent returns a new content.
func NewContent() *Content {
	c := Content{
		Pages: tview.NewPages(),
	}

	activatedCallback = make(map[string]func())
	c.Pages.SetChangedFunc(func() {
		if c.Pages.CurrentPage() != nil {
			if fn, ok := activatedCallback[c.Pages.CurrentPage().Name]; ok {
				fn()
			}
		}
	})

	return &c
}

func (c *Content) SubscribePageChange(viewId int, fn func()) {
	activatedCallback[views.Views[viewId].Name] = fn
}
