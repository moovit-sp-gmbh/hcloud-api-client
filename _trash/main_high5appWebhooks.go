package ui

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
)

type High5Webhook struct {
	hcloud.Webhook

	Selected bool
}

func (a *App) buildHigh5AppWebhooks(appId string) {
	table := a.getMainTable("helmut.cloud > service > high5 > app > webhooks")
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Event").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Target").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 3, tview.NewTableCell("Token").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 4, tview.NewTableCell("URL").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load high5 apps
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restWebhooks, err := high5.GetWebhooks(appId, 1000, 0)
	webhooks := []High5Webhook{}
	if err != nil {
		a.showError(err)
	} else {
		for _, webhook := range *restWebhooks {
			webhooks = append(webhooks, High5Webhook{Webhook: webhook, Selected: false})
		}

		for i, webhook := range webhooks {
			table.SetCell(i+2, 0, tview.NewTableCell(webhook.Name).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 1, tview.NewTableCell(webhook.EventId).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 2, tview.NewTableCell(webhook.Target).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 3, tview.NewTableCell(webhook.Token).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 4, tview.NewTableCell(webhook.Url).SetTextColor(tcell.ColorWhite))
		}
	}

	a.initLastSelection(table, main_high5_app_settings_webhooks, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			a.buildHigh5AppDetail(appId)
		default:
			a.buildHigh5AppWebhookLog(appId, webhooks[row-2].Id)
		}
	})

	table.registerKey(27, func() {
		a.buildHigh5AppDetail(appId)
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
