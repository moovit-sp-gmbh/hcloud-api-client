package views

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

func (m *Main) buildIdpOrganizationMembers(organizationId string) {
	table := m.getMainTable("helmut.cloud > service > idp > organization > members")

	table.SetCell(0, 0, tview.NewTableCell("E-Mail").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Role").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load idp organizations
	ctx := config.Config.GetActiveContext()
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restOrganizationMembers, err := idp.ListOrganizationMembersById(organizationId, 1000, 0)
	if err != nil {
		m.showError(err)
		return
	}

	organizationMembers := []IdpOrganizationMember{}
	for _, organizationMember := range *restOrganizationMembers {
		organizationMembers = append(organizationMembers, IdpOrganizationMember{OrganizationMember: organizationMember, Selected: false})
	}

	for i, organizationMember := range organizationMembers {
		color := tcell.ColorWhite
		if organizationMember.User.Id == m.ActiveUser.Id {
			color = tcell.ColorDarkOrange
		}
		table.SetCell(i+2, 0, tview.NewTableCell(organizationMember.User.Email).SetTextColor(color))
		table.SetCell(i+2, 1, tview.NewTableCell(organizationMember.User.Name).SetTextColor(tcell.ColorWhite))
		table.SetCell(i+2, 2, tview.NewTableCell(string(organizationMember.Permission)).SetTextColor(tcell.ColorWhite))
	}

	m.initLastSelection(table, main_idp_organizations, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildIdpOrganizations()
		default:
		}
	})

	table.registerKey(27, func() {
		m.buildIdpOrganizations()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
