package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Main) showConfirm(header string, width int, height int, callback func(bool), text ...string) {
	old := m.Main.Content.GetItem(0)
	form := tview.NewForm()
	g := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	for _, t := range text {
		form.AddTextView(t, t, 1, 1, false, false)
	}
	form.AddButton("Yes", func() {
		m.Main.Content.Clear().AddItem(old, 0, 1, false)
		m.SetFocus(old)
		callback(true)
	})
	form.AddButton("No", func() {
		m.Main.Content.Clear().AddItem(old, 0, 1, false)
		m.SetFocus(old)

		callback(false)
	})
	form.SetBorder(true).SetTitle("   " + header + "   ")

	m.Main.Content.Clear().AddItem(g, 0, 1, true)
	m.SetFocus(form)
}
