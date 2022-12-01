package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) showHelp() {
	table := a.getMainTable("helmut.cloud > help")
	table.SetCell(0, 0, tview.NewTableCell("shortcut").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("description").SetTextColor(tcell.ColorBlue))

	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("go one level up").SetTextColor(tcell.ColorWhite))

	table.SetCell(2, 0, tview.NewTableCell("enter").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(2, 1, tview.NewTableCell("step into selected entry").SetTextColor(tcell.ColorWhite))

	table.SetCell(3, 0, tview.NewTableCell("esc / arrow-left").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(3, 1, tview.NewTableCell("go one level up").SetTextColor(tcell.ColorWhite))

	table.SetCell(4, 0, tview.NewTableCell("ctrl+n").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(4, 1, tview.NewTableCell("create a new entity").SetTextColor(tcell.ColorWhite))

	table.SetCell(5, 0, tview.NewTableCell("ctrl+d").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(5, 1, tview.NewTableCell("delete an entity").SetTextColor(tcell.ColorWhite))

	table.SetCell(6, 0, tview.NewTableCell(":").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(6, 1, tview.NewTableCell("open command prompt (quick access)").SetTextColor(tcell.ColorWhite))

	table.SetCell(7, 0, tview.NewTableCell("   ctx / context").SetTextColor(tcell.ColorWhite))
	table.SetCell(7, 1, tview.NewTableCell("   open context").SetTextColor(tcell.ColorWhite))

	table.SetCell(8, 0, tview.NewTableCell("   svc / service").SetTextColor(tcell.ColorWhite))
	table.SetCell(8, 1, tview.NewTableCell("   open services").SetTextColor(tcell.ColorWhite))

	table.SetCell(9, 0, tview.NewTableCell("   home").SetTextColor(tcell.ColorWhite))
	table.SetCell(9, 1, tview.NewTableCell("   go to first level").SetTextColor(tcell.ColorWhite))

	table.SetCell(10, 0, tview.NewTableCell("   idp").SetTextColor(tcell.ColorWhite))
	table.SetCell(10, 1, tview.NewTableCell("   go to idp service").SetTextColor(tcell.ColorWhite))

	table.SetCell(11, 0, tview.NewTableCell("   high5").SetTextColor(tcell.ColorWhite))
	table.SetCell(11, 1, tview.NewTableCell("   go to high5 service").SetTextColor(tcell.ColorWhite))

	table.SetCell(12, 0, tview.NewTableCell("   help").SetTextColor(tcell.ColorWhite))
	table.SetCell(12, 1, tview.NewTableCell("   open this help page").SetTextColor(tcell.ColorWhite))

	table.SetCell(13, 0, tview.NewTableCell("   quit").SetTextColor(tcell.ColorWhite))
	table.SetCell(13, 1, tview.NewTableCell("   quit the program").SetTextColor(tcell.ColorWhite))

	table.SetCell(14, 0, tview.NewTableCell("?").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(14, 1, tview.NewTableCell("open help").SetTextColor(tcell.ColorWhite))

	table.SetCell(15, 0, tview.NewTableCell("ctrl+r").SetTextColor(tcell.ColorDarkOrange))
	table.SetCell(15, 1, tview.NewTableCell("reload current view").SetTextColor(tcell.ColorWhite))

	a.initLastSelection(table, main_help, 1)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 1:
			a.buildHcloud()
		}
	})

	table.registerKey(27, func() {
		a.buildHcloud()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)

}
