package views

import (
	"hcloud-api-client/config"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

func (m *Main) buildHigh5Webhooks(appId string) {
	table := m.getMainTable("helmut.cloud > service > high5 > app > webhooks")
	table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Event").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Target").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 3, tview.NewTableCell("Token").SetTextColor(tcell.ColorBlue))
	// table.SetCell(0, 4, tview.NewTableCell("URL").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 4, tview.NewTableCell("Security Headers").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load high5 apps
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restWebhooks, err := high5.GetWebhooks(appId, 1000, 0)
	webhooks := []High5Webhook{}
	if err != nil {
		m.showError(err)
	} else {
		for _, webhook := range *restWebhooks {
			webhooks = append(webhooks, High5Webhook{Webhook: webhook, Selected: false})
		}

		for i, webhook := range webhooks {
			headers := []string{}
			for k, v := range webhook.SecurityHeaders {
				headers = append(headers, k+": "+v)
			}
			table.SetCell(i+2, 0, tview.NewTableCell(webhook.Name).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 1, tview.NewTableCell(webhook.EventId).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 2, tview.NewTableCell(webhook.Target).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 3, tview.NewTableCell(webhook.Token).SetTextColor(tcell.ColorWhite))
			table.SetCell(i+2, 4, tview.NewTableCell(strings.Join(headers, ", ")).SetTextColor(tcell.ColorWhite))
		}
	}

	m.initLastSelection(table, main_high5_app_settings_webhooks, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildHigh5App(appId)
		default:
			m.buildHigh5WebhookLogs(appId, webhooks[row-2].Id)
		}
	})

	table.registerKey(27, func() {
		m.buildHigh5App(appId)
	})

	table.registerKey(int32(tcell.KeyCtrlD), func() {
		m.showConfirm("Please confirm", 50, 7, func(b bool) {
			if b {
				row, _ := table.GetSelection()
				err := high5.DeleteWebhookById(webhooks[row-2].Id)
				if err != nil {
					m.showError(err)
				} else {
					m.buildHigh5Webhooks(appId)
				}
			} else {
				m.buildHigh5Webhooks(appId)
			}
		}, "Are you sure you want to delete this webhook?")
	})

	table.registerKey(int32(tcell.KeyCtrlN), func() {
		var name, event, pat, target string

		restEvents, err := high5.GetEvents(appId, 1000, 0)
		if err != nil {
			m.showError(err)
			return
		}
		var events = []string{}
		for _, event := range *restEvents {
			events = append(events, event.Name)
		}

		restMembers, err := idp.ListOrganizationMembersById(m.ActiveOrganization.Id, 1000, 0)
		if err != nil {
			m.showError(err)
			return
		}
		var members = []string{}
		for _, member := range *restMembers {
			members = append(members, member.User.Email)
		}

		form, _ := m.newForm("Add new Webhook", 60, 13)
		form.AddFormItem(m.newInputField("Name", &name)).
			AddFormItem(m.newSelectField("Event", events, &event).SetFieldWidth(100)).
			AddFormItem(m.newInputField("PAT", &pat)).
			AddFormItem(m.newSelectField("Target", members, &target).SetFieldWidth(100)).
			AddButton("Save", func() {
				var eventId string
				for _, e := range *restEvents {
					if e.Name == event {
						eventId = e.Id
					}
				}

				_, err := high5.CreateWebhook(appId, eventId, name, pat, target, nil)
				if err != nil {
					m.showError(err)
				} else {
					m.buildHigh5Webhooks(appId)
				}
			}).
			AddButton("Cancel", func() {
				m.buildHigh5Webhooks(appId)
			})
		m.SetFocus(form)
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
