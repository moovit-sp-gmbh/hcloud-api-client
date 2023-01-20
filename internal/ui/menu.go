package ui

import (
	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type Menu struct {
	*tview.Table
}

// NewMenu returns a new menu.
func NewMenu() *Menu {
	m := Menu{
		Table: tview.NewTable(),
	}

	return &m
}

// BuildMenu populate menu ui KeyActions.
func (m *Menu) BuildMenu(ka KeyActions) *Menu {
	m.Clear()

	count := 0
	for _, k := range ka {
		if k.Visible {
			m.SetCell(count, 0, tview.NewTableCell("<"+k.Key+"> ").SetTextColor(tcell.ColorBlue))
			m.SetCell(count, 1, tview.NewTableCell(k.Description).SetTextColor(tcell.ColorWhite))
			count++
		}
	}

	return m
}
