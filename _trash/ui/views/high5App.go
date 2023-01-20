package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Main) buildHigh5App(appId string) {
	table := m.getMainTable("helmut.cloud > service > high5 > app")
	table.SetCell(0, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 0, tview.NewTableCell("events").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 0, tview.NewTableCell("webhooks").SetTextColor(tcell.ColorWhite))

	m.initLastSelection(table, main_high5_app_detail, 1)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			m.buildHigh5Apps()
		case 1:
			m.buildHigh5Events(appId)
		case 2:
			m.buildHigh5Webhooks(appId)
		}
	})

	table.registerKey(27, func() {
		m.buildHigh5Apps()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
