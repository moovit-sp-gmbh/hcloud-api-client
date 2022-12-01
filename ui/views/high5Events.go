package views

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
)

func (m *Main) buildHigh5Events(appId string) {
	table := m.getMainTable("helmut.cloud > service > high5 > app > events")
	table.SetCell(0, 0, tview.NewTableCell("Event").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load events for app
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restEvents, _ := high5.GetEvents(appId, 1000, 0)

	events := []High5Event{}
	for _, event := range *restEvents {
		events = append(events, High5Event{Event: event, Selected: false})
	}

	for i, event := range events {
		table.SetCell(i+2, 0, tview.NewTableCell(event.Name).SetTextColor(tcell.ColorWhite))
	}

	m.initLastSelection(table, main_high5_app_events, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildHigh5App(appId)
		default:
			m.buildHigh5Streams(appId, events[row-2].Id)
		}
	})

	table.registerKey(27, func() {
		m.buildHigh5App(appId)
	})
	table.registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		if row > 1 {
			m.showConfirm("Please confirm", 80, 7, func(b bool) {
				if b {
					err := high5.DeleteEventById(events[row-2].Id)
					if err == nil {
						m.buildHigh5Events(appId)
					} else {
						m.showError(err)
					}
				}
			}, "Are you sure you want to delete this event?")
		}
	})
	table.registerKey(int32(tcell.KeyCtrlN), func() {
		createHigh5Event(m, appId, high5)
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}

func createHigh5Event(m *Main, appId string, high5 *high5.Client) {
	name := ""
	form := tview.NewForm()
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	form.AddInputField("Event name", "", 20, nil, func(text string) { name = text }).
		AddButton("Save", func() {
			_, err := high5.CreateEvent(name, appId)
			if err != nil {
				m.showError(err)
			} else {
				m.buildHigh5Events(appId)
			}
		}).
		AddButton("Cancel", func() {
			m.buildHigh5Events(appId)
		})
	form.SetBorder(true).SetTitle("   Enter event name   ")

	g := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 7, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	m.Main.Content.Clear().AddItem(g, 0, 1, false)
	m.SetFocus(form)
}
