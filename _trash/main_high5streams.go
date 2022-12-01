package ui

import (
	"encoding/base64"
	"fmt"
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
	"github.com/skratchdot/open-golang/open"
)

type High5Stream struct {
	hcloud.Stream

	Selected bool
}

func (a *App) buildStreamContent(appId string, eventId string) {
	table := a.getMainTable("helmut.cloud > service > high5 > apps > events > streams")
	table.SetCell(0, 0, tview.NewTableCell("Stream").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load streams for event
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restStreams, _ := high5.GetStreams(eventId, 1000, 0)

	streams := []High5Stream{}
	for _, stream := range *restStreams {
		streams = append(streams, High5Stream{Stream: stream, Selected: false})
	}

	for i, stream := range streams {
		table.SetCell(i+2, 0, tview.NewTableCell(stream.Name).SetTextColor(tcell.ColorWhite))
	}

	a.initLastSelection(table, main_high5_app_event_streams, 2)
	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			a.buildEventContent(appId)
		default:
			url := fmt.Sprintf("%s/high5/designer/?b64jwt=%s#/%s", ctx.Server, base64.StdEncoding.EncodeToString([]byte(ctx.Token)), streams[row-2].Id)
			open.Run(url)
		}
	})

	table.registerKey(27, func() { a.buildEventContent(appId) })
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
			a.buildStreamContent(appId, eventId)
		} else {
			a.showError(err)
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
			a.buildStreamContent(appId, eventId)
		} else {
			a.showError(err)
		}
	})
	table.registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		if row > 1 {
			a.showConfirm("Please confirm", 80, 7, func(b bool) {
				if b {
					err := high5.DeleteStreamById(streams[row-2].Id)
					if err == nil {
						a.buildStreamContent(appId, eventId)
					} else {
						a.showError(err)
					}
				}
			}, "Are you sure you want to delete this stream?")
		}
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
					a.showError(err)
				} else {
					a.buildStreamContent(appId, eventId)
				}
			}).
			AddButton("Cancel", func() {
				a.buildStreamContent(appId, eventId)
			})
		form.SetBorder(true).SetTitle(addPaddingToString("Enter stream name"))

		g := tview.NewGrid().
			SetColumns(0, 80, 0).
			SetRows(0, 7, 0).
			AddItem(form, 1, 1, 1, 1, 0, 0, true)

		a.Pages.Main.Content.Clear().AddItem(g, 0, 1, false)
		a.SetFocus(form)
	})

	a.buildStreamsMenu()

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.SetFocus(table)
}

func (a *App) buildStreamsMenu() {
	table := tview.NewTable()
	table.SetCell(0, 0, tview.NewTableCell("<ctrl+n>").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("create new strean").SetTextColor(tcell.ColorWhite))

	table.SetCell(1, 0, tview.NewTableCell("<ctrl+d>").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 1, tview.NewTableCell("delete stream").SetTextColor(tcell.ColorWhite))

	table.SetCell(0, 2, tview.NewTableCell("<ctrl+h>").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 3, tview.NewTableCell("move stream up").SetTextColor(tcell.ColorWhite))

	table.SetCell(1, 2, tview.NewTableCell("<ctrl+l>").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 3, tview.NewTableCell("move stream down").SetTextColor(tcell.ColorWhite))

	a.Pages.Main.Menu.Clear().AddItem(table, 0, 1, false)
}
