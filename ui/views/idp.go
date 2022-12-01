package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/rivo/tview"
)

type IdpOrganization struct {
	hcloud.Organization

	Selected bool
}

type IdpOrganizationMember struct {
	hcloud.OrganizationMember

	Selected bool
}

func (m *Main) buildIdp() {
	table := m.getMainTable("helmut.cloud > service > idp")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("organizations").SetTextColor(tcell.ColorWhite))

	m.initLastSelection(table, main_high5_app_event_streams, 1)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			m.buildServices()
		case 1:
			m.buildIdpOrganizations()
		}
	})

	table.registerKey(27, func() {
		m.buildServices()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
