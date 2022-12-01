package views

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
)

func (m *Main) buildHigh5Apps() {
	table := m.getMainTable("helmut.cloud > service > high5 > apps")
	table.SetCell(0, 0, tview.NewTableCell("Context").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Creator").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Events").SetTextColor(tcell.ColorBlue))

	table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))

	// load high5 apps
	ctx := config.Config.GetActiveContext()
	high5 := high5.New(hcloud.New(&hcloud.Config{Api: ctx.Server, Token: ctx.Token}))
	restApps, _ := high5.GetApps(1000, 0)

	apps := []High5App{}
	for _, app := range *restApps {
		apps = append(apps, High5App{App: app, Selected: false})
	}

	for i, app := range apps {
		table.SetCell(i+2, 0, tview.NewTableCell(app.Name).SetTextColor(tcell.ColorWhite))
		table.SetCell(i+2, 1, tview.NewTableCell(app.CreatorId).SetTextColor(tcell.ColorWhite))
		table.SetCell(i+2, 2, tview.NewTableCell("0").SetTextColor(tcell.ColorWhite))
	}

	m.initLastSelection(table, main_high5_apps, 2)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			m.buildServices()
		default:
			m.buildHigh5App(apps[row-2].Id)
		}
	})

	table.registerKey(int32(tcell.KeyCtrlN), func() {
		createHigh5App(m, high5)
	})

	table.registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		if row > 1 {
			deleteHigh5App(m, high5, apps[row-2].Id)
		}
	})

	table.registerKey(27, func() {
		m.buildServices()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}

func deleteHigh5App(m *Main, high5 *high5.Client, appId string) {
	m.showConfirm("Please confirm", 45, 7, func(b bool) {
		if b {
			err := high5.DeleteAppById(appId)
			if err != nil {
				m.showError(err)
			} else {
				m.buildHigh5Apps()
			}
		}
	}, "Are you sure you want to delete this app?")
}

func createHigh5App(m *Main, high5 *high5.Client) {
	name := ""
	form := tview.NewForm()
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	form.AddInputField("App name", "", 20, nil, func(text string) { name = text }).
		AddButton("Save", func() {
			_, err := high5.CreateApp(name)
			if err != nil {
				m.showError(err)
			} else {
				m.buildHigh5Apps()
			}
		}).
		AddButton("Cancel", func() {
			m.buildHigh5Apps()
		})
	form.SetBorder(true).SetTitle("   Enter app name   ")

	g := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 7, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	m.Main.Content.Clear().AddItem(g, 0, 1, false)
	m.SetFocus(form)
}
