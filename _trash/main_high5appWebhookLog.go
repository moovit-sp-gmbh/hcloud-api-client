package ui

import (
	"fmt"
	"hcloud-api-client/config"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
)

type High5WebhookLog struct {
	hcloud.WebhookLog

	Selected bool
}

func (a *App) buildHigh5AppWebhookLog(appId string, webhookId string) {
	table := a.getMainTable("helmut.cloud > service > high5 > app > webhook > log")
	table.SetCell(0, 0, tview.NewTableCell("Timestamp").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Response Status").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Response Body").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load high5 apps
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restWebhookLogs, err := high5.GetWebhookLogs(webhookId, 1000, 0)
	if err != nil {
		a.showError(err)
	} else {
		webhookLogs := []High5WebhookLog{}
		for _, webhookLog := range *restWebhookLogs {
			webhookLogs = append(webhookLogs, High5WebhookLog{WebhookLog: webhookLog, Selected: false})
		}

		for i, webhookLog := range webhookLogs {
			tm := time.Unix(int64(webhookLog.Timestamp/1000), 0)
			table.SetCell(i+2, 0, tview.NewTableCell(fmt.Sprintf("%s", tm.String())).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 1, tview.NewTableCell(fmt.Sprintf("%d", webhookLog.ResponseStatusCode)).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 2, tview.NewTableCell(fmt.Sprintf("%s", webhookLog.ResponseBody)).SetTextColor(tcell.ColorWhite))
		}
	}

	a.initLastSelection(table, main_high5_app_settings_webhook_logs, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 1:
			a.buildHigh5AppWebhooks(appId)
		}
	})

	table.registerKey(27, func() {
		a.buildHigh5AppWebhooks(appId)
	})
	table.registerKey(int32(tcell.KeyCtrlR), func() {
		a.buildHigh5AppWebhookLog(appId, webhookId)
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
