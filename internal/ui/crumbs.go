package ui

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type BreadCrumb struct {
	Name  string
	Color tcell.Color
	Fn    func()
}

// Crumbs represents user breadcrumbs.
type Crumbs struct {
	*tview.Flex
}

// NewCrumbs returns a new breadcrumb view.
func NewCrumbs() *Crumbs {
	c := Crumbs{
		Flex: tview.NewFlex().SetDirection(tview.FlexColumn),
	}

	c.Flex.SetBorderPadding(1, 1, 1, 1)

	return &c
}

func (c *Crumbs) SetBreadcrumb(entries ...BreadCrumb) {
	c.Clear()

	for i, entry := range entries {
		if i > 0 {
			c.AddItem(tview.NewTextView().SetText(">").SetTextColor(tcell.ColorBlack), 2, 1, false)
		}
		tv := tview.NewTextView()
		tv.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			entry.Fn()
			return action, event
		})
		c.AddItem(tv.SetText(entry.Name).SetTextColor(entry.Color), len(entry.Name)+1, 1, false)
	}
}
