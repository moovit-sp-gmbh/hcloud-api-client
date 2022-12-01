package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) buildHigh5AppDetail(appId string) {
	table := a.getMainTable("helmut.cloud > service > high5 > app")
	table.SetCell(0, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 0, tview.NewTableCell("events").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 0, tview.NewTableCell("webhooks").SetTextColor(tcell.ColorWhite))

	a.initLastSelection(table, main_high5_app_detail, 1)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			a.buildHigh5Apps()
		case 1:
			a.buildEventContent(appId)
		case 2:
			a.buildHigh5AppWebhooks(appId)
		}
	})

	table.registerKey(27, func() {
		a.buildHigh5Apps()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
