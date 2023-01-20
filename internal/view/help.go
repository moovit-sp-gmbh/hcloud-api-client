package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type Help struct {
	*tableView
}

func NewHelp(a *ui.App) *Help {
	v := &Help{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.HELP,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *Help) setup() {
	v.createRow(0, nil, tableCell{name: "shortcut", color: tcell.ColorBlue}, tableCell{name: "shortcut", color: tcell.ColorBlue})

	v.createRow(1, func() {
		v.app.LoadContent(views.HCLOUD)
	}, tableCell{name: "..", color: tcell.ColorDarkOrange}, tableCell{name: "go one level up", color: tcell.ColorWhite})

	v.createRow(2, nil, tableCell{name: "enter", color: tcell.ColorDarkOrange}, tableCell{name: "step into selected entry", color: tcell.ColorWhite})

	v.createRow(3, nil, tableCell{name: "esc", color: tcell.ColorDarkOrange}, tableCell{name: "go one level up", color: tcell.ColorWhite})

	v.createRow(4, nil, tableCell{name: "ctrl+n", color: tcell.ColorDarkOrange}, tableCell{name: "create a new entity", color: tcell.ColorWhite})
	v.createRow(5, nil, tableCell{name: "ctrl+d", color: tcell.ColorDarkOrange}, tableCell{name: "delete selected entity", color: tcell.ColorWhite})
	v.createRow(6, nil, tableCell{name: "ctrl+r", color: tcell.ColorDarkOrange}, tableCell{name: "reload current view", color: tcell.ColorWhite})

	v.createRow(7, nil, tableCell{name: ":", color: tcell.ColorDarkOrange}, tableCell{name: "open command prompt", color: tcell.ColorWhite})
	v.createRow(8, nil, tableCell{name: "   ctx / context", color: tcell.ColorDarkOrange}, tableCell{name: "open context", color: tcell.ColorWhite})
	v.createRow(9, nil, tableCell{name: "   svc / service", color: tcell.ColorDarkOrange}, tableCell{name: "open service", color: tcell.ColorWhite})
	v.createRow(10, nil, tableCell{name: "   home", color: tcell.ColorDarkOrange}, tableCell{name: "go to first level", color: tcell.ColorWhite})
	v.createRow(11, nil, tableCell{name: "   idp", color: tcell.ColorDarkOrange}, tableCell{name: "go to idp service", color: tcell.ColorWhite})
	v.createRow(12, nil, tableCell{name: "   high5", color: tcell.ColorDarkOrange}, tableCell{name: "go to high5 service", color: tcell.ColorWhite})

	v.createRow(13, nil, tableCell{name: "   help", color: tcell.ColorDarkOrange}, tableCell{name: "open this help page", color: tcell.ColorWhite})
	v.createRow(14, nil, tableCell{name: "   quit", color: tcell.ColorDarkOrange}, tableCell{name: "quit the program", color: tcell.ColorWhite})

	v.createRow(15, nil, tableCell{name: "?", color: tcell.ColorDarkOrange}, tableCell{name: "open this help", color: tcell.ColorWhite})
	v.createRow(16, nil, tableCell{name: "ctrl+c", color: tcell.ColorDarkOrange}, tableCell{name: "quit the program", color: tcell.ColorWhite})
}

func (v *Help) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HELP].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.HELP) }})
}

func (v *Help) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "back to hcloud", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.HCLOUD)
		return ek
	}, true)
}

func (v *Help) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *Help) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *Help) pageLoaded() {}
