package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) buildServices() {
	table := a.getMainTable("helmut.cloud > service")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("idp").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 1, tview.NewTableCell("high5").SetTextColor(tcell.ColorWhite))
	table.SetCell(3, 1, tview.NewTableCell("fuse").SetTextColor(tcell.ColorGray))
	table.SetCell(4, 1, tview.NewTableCell("auditor").SetTextColor(tcell.ColorGray))
	table.SetCell(5, 1, tview.NewTableCell("mailer").SetTextColor(tcell.ColorGray))

	a.initLastSelection(table, main_services, 1)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			a.buildHcloud()
		case 1:
			a.buildIdp()
		case 2:
			a.buildHigh5Apps()
		}
	})

	table.registerKey(27, func() {
		a.buildHcloud()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
