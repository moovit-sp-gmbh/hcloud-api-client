package views

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

func (m *Main) buildIdpOrganizations() {
	table := m.getMainTable("helmut.cloud > service > idp > organizations")

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
		color := tcell.ColorWhite
		if organization.Id == m.ActiveOrganization.Id {
			color = tcell.ColorDarkOrange
		}
		table.SetCell(i+2, 0, tview.NewTableCell(organization.Name).SetTextColor(color))
		table.SetCell(i+2, 1, tview.NewTableCell(organization.Company).SetTextColor(color))
	}

	m.initLastSelection(table, main_idp_organizations, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildIdp()
		default:
			m.buildIdpOrganizationMembers(organizations[row-2].Id)
		}
	})

	table.registerKey(int32(tcell.KeyCtrlA), func() {
		row, _ := table.GetSelection()
		if m.ActiveOrganization.Id != organizations[row-2].Id {
			user, err := idp.PatchUser(hcloud.PatchUser{ActiveOrganizationId: organizations[row-2].Id})
			if err != nil {
				m.showError(err)
			} else {
				m.ActiveUser = user
				organization, err := idp.GetOrganizationById(m.ActiveUser.ActiveOrganizationId)
				if err != nil {
					m.showError(err)
				} else {
					m.ActiveOrganization = organization
				}
				m.buildIdpOrganizations()
				m.buildInfo()
			}
		}
	})

	table.registerKey(27, func() {
		m.buildIdp()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
