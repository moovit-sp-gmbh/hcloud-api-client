package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type IDP struct {
	*tableView
}

func NewIDP(a *ui.App) *IDP {
	v := &IDP{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.SERVICE_IDP,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *IDP) setup() {
	v.createRow(0, func() { v.app.LoadContent(views.SERVICE) }, tableCell{name: "..", color: tcell.ColorWhite})
	v.createRow(1, func() { v.app.LoadContent(views.SERVICE_IDP_ACCOUNT) }, tableCell{name: "account", color: tcell.ColorWhite})
	v.createRow(2, func() { v.app.LoadContent(views.SERVICE_IDP_ORGANIZATIONS) }, tableCell{name: "organizations", color: tcell.ColorWhite})
}

func (v *IDP) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.SERVICE].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.SERVICE) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.SERVICE_IDP].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.SERVICE_IDP) }})
}

func (v *IDP) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "back to service", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.SERVICE)
		return ek
	}, true)
}

func (v *IDP) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *IDP) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *IDP) pageLoaded() {}
