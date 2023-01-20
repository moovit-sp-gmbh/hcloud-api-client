package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// Hcloud is the entrypoint for the tui
type Hcloud struct {
	*tableView
}

// NewHcloud returns a fully setup nwe Hcloud struct
func NewHcloud(a *ui.App) *Hcloud {
	v := &Hcloud{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.HCLOUD,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 0,
		},
	}

	v.setupView(v)

	return v
}

func (v *Hcloud) setup() {
	v.createRow(0, func() { v.app.LoadContent(views.SERVICE) }, tableCell{name: "service", color: tcell.ColorWhite})
	v.createRow(1, func() { v.app.LoadContent(views.CONFIG) }, tableCell{name: "config", color: tcell.ColorWhite})
	v.createRow(2, func() { v.app.LoadContent(views.HELP) }, tableCell{name: "help", color: tcell.ColorWhite})

	if v.app.IsDebugEnabled() {
		v.createRow(3, func() { v.app.LoadContent(views.DEBUG) }, tableCell{name: "debug", color: tcell.ColorWhite})
	}

}

func (v *Hcloud) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
}

func (v *Hcloud) setupKeyActions() {
	// no key actions for this view
}

func (v *Hcloud) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *Hcloud) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *Hcloud) pageLoaded() {}
