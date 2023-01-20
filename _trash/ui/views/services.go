package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Main) buildServices() {
	table := m.getMainTable("helmut.cloud > service")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("idp").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 1, tview.NewTableCell("high5").SetTextColor(tcell.ColorWhite))
	table.SetCell(3, 1, tview.NewTableCell("fuse").SetTextColor(tcell.ColorGray))
	table.SetCell(4, 1, tview.NewTableCell("auditor").SetTextColor(tcell.ColorGray))
	table.SetCell(5, 1, tview.NewTableCell("mailer").SetTextColor(tcell.ColorGray))

	m.initLastSelection(table, main_services, 1)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			m.BuildHcloud()
		case 1:
			m.buildIdp()
		case 2:
			m.buildHigh5()
		}
	})

	table.registerKey(27, func() {
		m.BuildHcloud()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
