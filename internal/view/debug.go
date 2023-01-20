package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"
	"time"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type Debug struct {
	*tableView
}

func NewDebug(a *ui.App) *Debug {
	v := &Debug{
		tableView: &tableView{
			Table: tview.NewTable(),

			app:  a,
			view: views.DEBUG,

			keyActions: make(ui.KeyActions),
			breadCrumb: []ui.BreadCrumb{},
			selectedFn: make(map[int]func()),

			selectedRow: 1,
		},
	}

	v.setupView(v)

	return v
}

func (v *Debug) setup() {
	v.Clear()
	v.createRow(0, nil, tableCell{name: "Time", color: tcell.ColorBlue}, tableCell{name: "Message", color: tcell.ColorBlue})

	v.createRow(1, func() {
		v.app.LoadContent(views.HCLOUD)
	}, tableCell{name: "..", color: tcell.ColorWhite})

	for i, debugLog := range v.app.GetDebugLogs() {
		v.createRow(i+2, nil, tableCell{name: time.Unix(debugLog.Timestamp/1000/1000/1000, 0).Format("2006-01-02 15:04:05"), color: tcell.ColorWhite}, tableCell{name: debugLog.Message, color: tcell.ColorWhite})
	}

	v.Select(v.GetRowCount()-1, 0)

}

func (v *Debug) setupBreadcrumb() {
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.HCLOUD].HumanReadableName, Color: tcell.ColorYellow, Fn: func() { v.app.LoadContent(views.HCLOUD) }})
	v.breadCrumb = append(v.breadCrumb, ui.BreadCrumb{Name: views.Views[views.DEBUG].HumanReadableName, Color: tcell.ColorDarkOrange, Fn: func() { v.app.LoadContent(views.DEBUG) }})
}

func (v *Debug) setupKeyActions() {
	v.keyActions[tcell.KeyEsc] = ui.NewKeyAction("esc", "back to hcloud", func(ek *tcell.EventKey) *tcell.EventKey {
		v.app.LoadContent(views.HCLOUD)
		return ek
	}, true)
}

func (v *Debug) drawMenu() {
	v.app.Menu().BuildMenu(v.keyActions)
}

func (v *Debug) drawBreadcrumb() {
	v.app.Crumbs().SetBreadcrumb(v.breadCrumb...)
}

func (v *Debug) pageLoaded() {
	v.setup()
}
