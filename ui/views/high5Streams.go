package views

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
	"github.com/skratchdot/open-golang/open"
)

func (m *Main) buildHigh5Streams(appId string, eventId string) {
	table := m.getMainTable("helmut.cloud > service > high5 > apps > events > streams")
	table.SetCell(0, 0, tview.NewTableCell("Stream").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load streams for event
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	idp := idp.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restStreams, _ := high5.GetStreams(eventId, 1000, 0)

	streams := []High5Stream{}
	for _, stream := range *restStreams {
		streams = append(streams, High5Stream{Stream: stream, Selected: false})
	}

	for i, stream := range streams {
		table.SetCell(i+2, 0, tview.NewTableCell(stream.Name).SetTextColor(tcell.ColorWhite))
	}

	m.initLastSelection(table, main_high5_app_event_streams, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildHigh5Events(appId)
		default:
			url := fmt.Sprintf("%s/high5/designer/?b64jwt=%s#/%s", ctx.Server, base64.StdEncoding.EncodeToString([]byte(ctx.Token)), streams[row-2].Id)
			open.Run(url)
		}
	})

	table.registerKey(27, func() { m.buildHigh5Events(appId) })
	table.registerKey(int32(tcell.KeyCtrlH), func() {
		streamOrder := []hcloud.StreamOrder{}
		row, _ := table.GetSelection()
		for _, stream := range streams {
			if stream.Id == streams[row-2].Id {
				if stream.Order > 1 {
					for _, s := range streams {
						if s.Order == stream.Order-1 {
							s.Order += 1
						}
					}
					stream.Order -= 1
				}
			}
		}
		for _, stream := range streams {
			streamOrder = append(streamOrder, hcloud.StreamOrder{StreamId: stream.Id, Order: stream.Order})
		}
		_, err := high5.ChangeStreamOrder(eventId, streamOrder)
		if err == nil {
			m.buildHigh5Streams(appId, eventId)
		} else {
			m.showError(err)
		}
	})
	table.registerKey(int32(tcell.KeyCtrlL), func() {
		streamOrder := []hcloud.StreamOrder{}
		row, _ := table.GetSelection()
		for _, stream := range streams {
			if stream.Id == streams[row-2].Id {
				if stream.Order < len(streams) {
					for _, s := range streams {
						if s.Order == stream.Order+1 {
							s.Order -= 1
						}
					}
					stream.Order += 1
				}
			}
		}
		for _, stream := range streams {
			streamOrder = append(streamOrder, hcloud.StreamOrder{StreamId: stream.Id, Order: stream.Order})
		}
		_, err := high5.ChangeStreamOrder(eventId, streamOrder)
		if err == nil {
			m.buildHigh5Streams(appId, eventId)
		} else {
			m.showError(err)
		}
	})
	table.registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		if row > 1 {
			m.showConfirm("Please confirm", 80, 7, func(b bool) {
				if b {
					err := high5.DeleteStreamById(streams[row-2].Id)
					if err == nil {
						m.buildHigh5Streams(appId, eventId)
					} else {
						m.showError(err)
					}
				}
			}, "Are you sure you want to delete this stream?")
		}
	})

	table.registerKey(int32(tcell.KeyCtrlE), func() {
		restMembers, err := idp.ListOrganizationMembersById(m.ActiveOrganization.Id, 1000, 0)
		if err != nil {
			m.showError(err)
			return
		}
		var members = []string{}
		for _, member := range *restMembers {
			members = append(members, member.User.Email)
		}

		var target, dataType, payload string

		var form *tview.Form
		var formGrid *tview.Grid
		form, formGrid = m.newForm("Execute stream", 100, 60)
		form.AddFormItem(m.newSelectField("Target", members, &target)).
			AddFormItem(m.newSelectField("Payload Type", []string{"JSON", "GENERIC"}, &dataType)).
			AddFormItem(m.newTextArea("Payload Type", 100, 5, &payload)).
			AddButton("Execute", func() {
				row, _ := table.GetSelection()

				res, err := high5.ExecuteStreamById(streams[row-2].Id, target, dataType, []byte(payload), 10000, true)
				if err != nil {
					m.showError(err)
				} else {
					byteRes, _ := json.Marshal(res)
					res := string(byteRes)
					respForm, _ := m.newForm("Execution result", 100, 100)
					respForm.AddFormItem(m.newTextArea("Result", 100, 16, &res)).
						AddButton("Close", func() {
							m.Main.Content.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
							m.Main.Content.Clear().AddItem(formGrid, 0, 1, false)
							m.SetFocus(form)
						})
					m.SetFocus(respForm)
				}

			}).
			AddButton("Cancel", func() {
				m.buildHigh5Streams(appId, eventId)
			})

		m.SetFocus(form)
	})
	table.registerKey(int32(tcell.KeyCtrlN), func() {
		name := ""
		form := tview.NewForm()
		form.SetFieldTextColor(tcell.ColorWhite)
		form.SetFieldBackgroundColor(tcell.ColorBlue)
		form.AddInputField("Stream name", "", 20, nil, func(text string) { name = text }).
			AddButton("Save", func() {
				_, err := high5.CreateStream(name, eventId)
				if err != nil {
					m.showError(err)
				} else {
					m.buildHigh5Streams(appId, eventId)
				}
			}).
			AddButton("Cancel", func() {
				m.buildHigh5Streams(appId, eventId)
			})
		form.SetBorder(true).SetTitle("   Enter stream name   ")

		g := tview.NewGrid().
			SetColumns(0, 80, 0).
			SetRows(0, 7, 0).
			AddItem(form, 1, 1, 1, 1, 0, 0, true)

		m.Main.Content.Clear().AddItem(g, 0, 1, false)
		m.SetFocus(form)
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.SetFocus(table)
}
