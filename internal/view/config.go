package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"
	"hcloud-api-client/pkg"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/skratchdot/open-golang/open"
)

type Config struct {
	*tableView
}

func NewConfig(a *ui.App) *Config {
	v := &Config{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.CONFIG,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *Config) setup() {
	v.createRow(0, func() { v.app.LoadContent(views.HCLOUD) }, tableCell{name: "..", color: tcell.ColorWhite})
	v.createRow(1, func() { v.app.LoadContent(views.CONFIG_CONTEXT) }, tableCell{name: "context", color: tcell.ColorWhite})
	v.createRow(2, func() { open.Run(pkg.GetHomeDir() + "/.hcloud/") }, tableCell{name: "open config folder", color: tcell.ColorWhite})
	v.createRow(3, func() { open.Run(pkg.GetHomeDir() + "/.hcloud/config.yml") }, tableCell{name: "open config file", color: tcell.ColorWhite})
}

func (v *Config) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.CONFIG].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.CONFIG) }})
}

func (v *Config) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "YO back to hcloud", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.HCLOUD)
		return ek
	}, true)
}

func (v *Config) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *Config) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *Config) pageLoaded() {}
