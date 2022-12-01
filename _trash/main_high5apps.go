package ui

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/high5"
	"github.com/rivo/tview"
)

type High5App struct {
	hcloud.App

	Selected bool
}

func (a *App) buildHigh5Apps() {
	a.buildAppContent()
	a.buildMenu()
}

func (a *App) buildMenu() {
	table := tview.NewTable()
	table.SetCell(0, 0, tview.NewTableCell("<ctrl+n>").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("create new app").SetTextColor(tcell.ColorWhite))

	table.SetCell(1, 0, tview.NewTableCell("<ctrl+d>").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 1, tview.NewTableCell("delete app").SetTextColor(tcell.ColorWhite))

	a.Pages.Main.Menu.Clear().AddItem(table, 0, 1, false)
}

func (a *App) buildAppContent() {
	table := a.getMainTable("helmut.cloud > service > high5 > apps")
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

	a.initLastSelection(table, main_high5_apps, 2)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			a.buildServices()
		default:
			a.buildHigh5AppDetail(apps[row-2].Id)
		}
	})

	table.registerKey(int32(tcell.KeyCtrlN), func() {
		a.createNewApp(high5)
	})

	table.registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		if row > 1 {
			a.deleteHigh5App(high5, apps[row-2].Id)
		}
	})

	table.registerKey(27, func() {
		a.buildServices()
	})

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}

func (a *App) deleteHigh5App(high5 *high5.Client, appId string) {
	a.showConfirm("Please confirm", 80, 11, func(b bool) {
		if b {
			err := high5.DeleteAppById(appId)
			if err != nil {
				a.showError(err)
			} else {
				a.buildHigh5Apps()
			}
		}
	}, "Are you sure you want to delete this app?")
}

func (a *App) createNewApp(high5 *high5.Client) {
	name := ""
	form := tview.NewForm()
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	form.AddInputField("App name", "", 20, nil, func(text string) { name = text }).
		AddButton("Save", func() {
			_, err := high5.CreateApp(name)
			if err != nil {
				a.showError(err)
			} else {
				a.buildHigh5Apps()
			}
		}).
		AddButton("Cancel", func() {
			a.buildHigh5Apps()
		})
	form.SetBorder(true).SetTitle(addPaddingToString("Enter app name"))

	g := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 7, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	a.Pages.Main.Content.Clear().AddItem(g, 0, 1, false)
	a.SetFocus(form)
}
