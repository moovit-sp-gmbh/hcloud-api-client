package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type High5 struct {
	*tableView
}

func NewHigh5(a *ui.App) *High5 {
	v := &High5{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.SERVICE_HIGH5,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *High5) setup() {
	v.createRow(0, func() { v.app.LoadContent(views.SERVICE) }, tableCell{name: "..", color: tcell.ColorWhite})
	v.createRow(1, func() { v.app.LoadContent(views.SERVICE_HIGH5) }, tableCell{name: "apps", color: tcell.ColorWhite})
}

func (v *High5) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.SERVICE].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.SERVICE) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.SERVICE_HIGH5].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.SERVICE_HIGH5) }})
}

func (v *High5) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "back to service", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.SERVICE)
		return ek
	}, true)
}

func (v *High5) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *High5) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *High5) pageLoaded() {}
