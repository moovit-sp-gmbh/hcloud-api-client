package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) buildIdp() {
	table := a.getMainTable("helmut.cloud > service > idp")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("organizations").SetTextColor(tcell.ColorWhite))

	a.initLastSelection(table, main_high5_app_event_streams, 1)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			a.buildServices()
		case 1:
			a.buildIdpOrganizations()
		}
	})

	table.registerKey(27, func() {
		a.buildServices()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
