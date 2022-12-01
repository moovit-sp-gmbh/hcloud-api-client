package ui

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

func (a *App) buildHcloud() {
	table := a.getMainTable("helmut.cloud")

	table.SetCell(0, 1, tview.NewTableCell("service").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("config").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 1, tview.NewTableCell("help").SetTextColor(tcell.ColorWhite))

	ctx := config.Config.Active
	a.buildInfo()

	if ctx != nil {
		idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
		user, err := idp.Authorize()
		if err != nil {
			a.buildConfigContext(true)
			return
		} else {
			activeUser = user
			activeOrganization, _ = idp.GetOrganizationById(activeUser.ActiveOrganizationId)
			a.buildInfo()
		}
	} else {
		a.buildConfigContext(true)
		return
	}

	a.initLastSelection(table, main_hcloud, 0)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			if activeUser != nil {
				a.buildServices()
			}
		case 1:
			a.buildConfig()
		case 2:
			a.showHelp()
		}
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.SetFocus(table)
}
