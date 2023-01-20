package ui

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type Context struct {
	*tview.Table
}

// NewContext returns a new menu.
func NewContext() *Context {
	c := Context{
		Table: tview.NewTable(),
	}

	c.SetBorderPadding(0, 0, 1, 0)

	c.SetCell(0, 0, tview.NewTableCell("Context:").SetTextColor(tcell.ColorDarkOrange))
	c.SetCell(0, 1, tview.NewTableCell("dev").SetTextColor(tcell.ColorWhite))

	c.SetCell(1, 0, tview.NewTableCell("Server:").SetTextColor(tcell.ColorDarkOrange))
	c.SetCell(1, 1, tview.NewTableCell("https://dev.app.helmut.cloud").SetTextColor(tcell.ColorWhite))

	c.SetCell(2, 0, tview.NewTableCell("Account:").SetTextColor(tcell.ColorDarkOrange))
	c.SetCell(2, 1, tview.NewTableCell("account@moovit-sp.com").SetTextColor(tcell.ColorWhite))

	c.SetCell(3, 0, tview.NewTableCell("Organization:").SetTextColor(tcell.ColorYellow))
	c.SetCell(3, 1, tview.NewTableCell("account@moovit-sp.com").SetTextColor(tcell.ColorWhite))

	return &c
}
