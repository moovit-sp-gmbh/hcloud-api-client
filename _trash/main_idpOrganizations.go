package ui

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

type IdpOrganization struct {
	hcloud.Organization

	Selected bool
}

func (a *App) buildIdpOrganizations() {
	table := a.getMainTable("helmut.cloud > service > idp > organizations")

	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Company").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load idp organizations
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restOrganizations, _ := idp.ListOrganizations(1000, 0)

	organizations := []IdpOrganization{}
	for _, organization := range *restOrganizations {
		organizations = append(organizations, IdpOrganization{Organization: organization, Selected: false})
	}

	for i, organization := range organizations {
		table.SetCell(i+2, 0, tview.NewTableCell(organization.Name).SetTextColor(tcell.ColorWhite))
		table.SetCell(i+2, 1, tview.NewTableCell(organization.Company).SetTextColor(tcell.ColorWhite))
	}

	a.initLastSelection(table, main_idp_organizations, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			a.buildIdp()
		default:
		}
	})

	table.registerKey(27, func() {
		a.buildIdp()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
