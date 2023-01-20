package views

import "github.com/rivo/tview"

func (m *Main) resetMenu() {
	m.Main.Menu.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
}
