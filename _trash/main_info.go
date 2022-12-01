package ui

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) buildInfo() {
	table := tview.NewTable()

	ctx := config.Config.Active
	if ctx != nil && activeUser != nil && activeOrganization != nil {
		table.SetCell(0, 0, tview.NewTableCell("Context:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(0, 1, tview.NewTableCell(ctx.Identifier).SetTextColor(tcell.ColorWhite))
		table.SetCell(1, 0, tview.NewTableCell("Server:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(1, 1, tview.NewTableCell(ctx.Server).SetTextColor(tcell.ColorWhite))
		table.SetCell(2, 0, tview.NewTableCell("Account:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(2, 1, tview.NewTableCell(ctx.Email).SetTextColor(tcell.ColorWhite))
		table.SetCell(3, 0, tview.NewTableCell("Organization:").SetTextColor(tcell.ColorYellow))
		table.SetCell(3, 1, tview.NewTableCell(activeOrganization.Name).SetTextColor(tcell.ColorWhite))
	} else {
		table.SetCell(0, 0, tview.NewTableCell("Context:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(0, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorWhite))
		table.SetCell(1, 0, tview.NewTableCell("Server:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(1, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorWhite))
		table.SetCell(2, 0, tview.NewTableCell("Account:").SetTextColor(tcell.ColorOrangeRed))
		table.SetCell(2, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorWhite))
		table.SetCell(3, 0, tview.NewTableCell("Organization:").SetTextColor(tcell.ColorYellow))
		table.SetCell(3, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorWhite))
	}

	a.Pages.Main.Info.Clear().AddItem(table, 0, 1, false)
}
