package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// skipped for now as apps is the only option
func (a *App) buildHigh5() {
	table := a.getMainTable("helmut.cloud > service > high5")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("apps").SetTextColor(tcell.ColorWhite))

	a.initLastSelection(table, main_high5, 1)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			a.buildServices()
		case 1:
			a.buildHigh5Apps()
		}
	})

	table.registerKey(27, func() {
		a.buildServices()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
