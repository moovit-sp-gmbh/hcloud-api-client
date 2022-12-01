package ui

import (
	"hcloud-api-client/config"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go/service/idp"
	"github.com/rivo/tview"
)

func (a *App) buildConfigContext(lockUser bool) {
	table := a.getMainTable("helmut.cloud > config > context")

	table.SetCell(0, 0, tview.NewTableCell("Identifier").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("Server").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 2, tview.NewTableCell("Account").SetTextColor(tcell.ColorBlue))
	if !lockUser {
		table.SetCell(1, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	} else {
		table.SetCell(1, 0, tview.NewTableCell("").SetTextColor(tcell.ColorWhite))
	}

	defaultSelection := 0

	ctxs := config.Config.Contexts
	for i, ctx := range ctxs {
		color := tcell.ColorWhite
		if config.Config.Active != nil && config.Config.Active.Identifier == ctx.Identifier {
			color = tcell.ColorDarkOrange
			defaultSelection = i + 2
		}
		table.SetCell(i+2, 0, tview.NewTableCell(ctx.Identifier).SetTextColor(color))
		table.SetCell(i+2, 1, tview.NewTableCell(ctx.Server).SetTextColor(color))
		table.SetCell(i+2, 2, tview.NewTableCell(ctx.Email).SetTextColor(color))
	}

	a.initLastSelection(table, main_config_context, defaultSelection)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
		case 1:
			if !lockUser {
				a.buildConfig()
			}
		default:
			config.SetContext(ctxs[row-2].Identifier)
			// we clear the last selection history if we change context
			a.Pages.Main.SelectedRow = make(map[int]int)
			a.buildHcloud()
		}
	})

	table.registerKey(27, func() {
		a.buildConfig()
	}).registerKey(int32(tcell.KeyCtrlN), func() {
		identifier := ""
		server := ""
		email := ""
		password := ""
		totp := ""

		form := tview.NewForm().
			AddInputField("Identifier", "", 20, nil, func(text string) {
				identifier = text
			}).
			AddInputField("Server *", "", 20, nil, func(text string) {
				server = text
			}).
			AddInputField("E-Mail *", "", 20, nil, func(text string) {
				email = text
			}).
			AddPasswordField("Password *", "", 10, '*', func(text string) {
				password = text
			}).
			AddInputField("TOTP", "", 20, nil, func(text string) {
				totp = text
			}).
			AddButton("Save", func() {
				idp := idp.New(hcloud.New(&hcloud.Config{Api: server}))
				token, err := idp.Authenticate(email, password)
				if err != nil {

				} else {
					if identifier == "" {
						identifier = server
					}
					strings.Split(totp, "")
					config.AddContext(identifier, server, email, token.Token)
					config.SetContext(identifier)
					a.buildHcloud()
				}
			}).
			AddButton("Cancel", func() {
			})
		form.SetBorder(true).SetTitle("Create new context").SetTitleAlign(tview.AlignCenter)

		g := tview.NewGrid().
			SetColumns(0, 150, 0).
			SetRows(0, 80, 0).
			AddItem(form, 1, 1, 1, 1, 0, 0, true)

		a.Pages.Main.Content.Clear().AddItem(g, 0, 1, false)
		a.SetFocus(form)
	}).registerKey(int32(tcell.KeyCtrlD), func() {
		row, _ := table.GetSelection()
		config.DelContext(config.Config.Contexts[row-2].Identifier)
		table.RemoveRow(row)
		table.Select(2, 0)
	})

	a.buildConfigContextMenu()

	a.Pages.Main.Content.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}

func (a *App) buildConfigContextMenu() {
	table := getMenuTable()

	table.SetCell(0, 0, tview.NewTableCell("<ctrl+n>").SetTextColor(tcell.ColorBlue))
	table.SetCell(0, 1, tview.NewTableCell("new").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 0, tview.NewTableCell("<ctrl+d>").SetTextColor(tcell.ColorBlue))
	table.SetCell(1, 1, tview.NewTableCell("delete").SetTextColor(tcell.ColorWhite))
	a.Pages.Main.Menu.Clear().AddItem(table, 0, 1, false)
	a.Application.SetFocus(table)
}
