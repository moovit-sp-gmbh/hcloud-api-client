package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/rivo/tview"
)

type High5App struct {
	hcloud.App

	Selected bool
}

type High5Event struct {
	hcloud.Event

	Selected bool
}

type High5Webhook struct {
	hcloud.Webhook

	Selected bool
}

type High5WebhookLog struct {
	hcloud.WebhookLog

	Selected bool
}

type High5Stream struct {
	hcloud.Stream

	Selected bool
}

func (m *Main) buildHigh5() {
	table := m.getMainTable("helmut.cloud > service > high5")

	table.SetCell(0, 1, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 1, tview.NewTableCell("apps").SetTextColor(tcell.ColorWhite))

	m.initLastSelection(table, main_high5_app_event_streams, 1)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			m.buildServices()
		case 1:
			m.buildHigh5Apps()
		}
	})

	table.registerKey(27, func() {
		m.buildServices()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
