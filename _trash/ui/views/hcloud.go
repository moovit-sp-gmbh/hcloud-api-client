package views

import (
	"hcloud-api-client/config"
	"hcloud-api-client/ui/pages"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

type Main struct {
	*tview.Application
	*pages.Main
}

func (m *Main) BuildHcloud() {
	m.buildInfo()
	table := m.getMainTable("helmut.cloud")

	table.SetCell(0, 1, tview.NewTableCell("service").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("config").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 1, tview.NewTableCell("help").SetTextColor(tcell.ColorWhite))

	ctx := config.Config.Active

	if ctx != nil {
		idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
		user, err := idp.Authorize()
		if err != nil {
			m.buildContext(true)
			return
		} else {
			m.ActiveUser = user
			organization, err := idp.GetOrganizationById(m.ActiveUser.ActiveOrganizationId)
			if err != nil {
				m.showError(err)
			} else {
				m.ActiveOrganization = organization
			}

			m.buildInfo()
		}
	} else {
		m.buildContext(true)
		return
	}

	m.initLastSelection(table, main_hcloud, 0)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			if m.ActiveUser != nil {
				m.buildServices()
			}
		case 1:
			m.buildConfig()
		case 2:
			m.buildHelp()
		}
	})

	m.Content.Clear().AddItem(table, 0, 1, false)
	m.SetFocus(table)
}
