package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// Service is the entrypoint for the service tui
type Service struct {
	*tableView
}

// NewService returns a fully setup nwe Service struct
func NewService(a *ui.App) *Service {
	v := &Service{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.SERVICE,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *Service) setup() {
	v.createRow(0, func() { v.app.LoadContent(views.HCLOUD) }, tableCell{name: "..", color: tcell.ColorWhite})
	v.createRow(1, func() { v.app.LoadContent(views.SERVICE_IDP) }, tableCell{name: "idp", color: tcell.ColorWhite})
	v.createRow(2, func() { v.app.LoadContent(views.SERVICE_HIGH5) }, tableCell{name: "high5", color: tcell.ColorWhite})
}

func (v *Service) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.SERVICE].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.SERVICE) }})
}

func (v *Service) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "back to hcloud", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.HCLOUD)
		return ek
	}, true)
}

func (v *Service) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *Service) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *Service) pageLoaded() {}
