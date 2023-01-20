package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type TableCell struct {
	*tview.TableCell
}
type tableViewI interface {
	setup()

	setupBreadcrumb()
	setupKeyActions()

	drawBreadcrumb()
	drawMenu()

	pageLoaded()
}

// tableView represents a table-page within a.Content.
type tableView struct {
	*tview.Table
	app *ui.App

	view int

	keyActions ui.KeyActions
	breadCrumb []ui.BreadCrumb
	selectedFn map[int]func()

	selectedRow int
}

// setupView sets up a new view.
func (v *tableView) setupView(view tableViewI) {
	v.SetBorder(true)
	v.SetBorderColor(tcell.ColorGray)
	v.SetBorderFocusColor(tcell.ColorBlue)

	v.SetSelectable(true, false)
	v.SetTitle("   " + views.Views[v.view].HumanReadableName + "   ")

	v.Select(v.selectedRow, 0)
	v.SetSelectedFunc(func(row, column int) {
		if fn, ok := v.selectedFn[row]; ok {
			fn()
		}
	})

	v.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if keyAction, ok := v.keyActions[event.Key()]; ok {
			keyAction.Action(event)
		}
		return event
	})

	view.setup()
	v.app.Debug("setup %s for view name '%s' with id '%d'", "tableData", views.Views[v.view].Name, v.view)

	view.setupKeyActions()
	v.app.Debug("setup %s for view name '%s' with id '%d'", "keyActions", views.Views[v.view].Name, v.view)

	view.setupBreadcrumb()
	v.app.Debug("setup %s for view name '%s' with id '%d'", "breadCrumb", views.Views[v.view].Name, v.view)

	v.app.Content().SubscribePageChange(v.view, func() {
		view.drawMenu()
		view.drawBreadcrumb()

		view.pageLoaded()

		v.app.SetFocus(v)
	})
	v.app.Debug("subscribed to %s event for view name '%s' with id '%d'", "pageLoad", views.Views[v.view].Name, v.view)

	v.app.Debug("setup view '%s' with id '%d'", views.Views[v.view].Name, v.view)
	v.app.Content().AddPage(views.Views[v.view].Name, v, true, false)
}

func (v *tableView) registerKey(action ui.KeyAction) {

}

type tableCell struct {
	name  string
	color tcell.Color
}

// createRow creates new cells for a complete road using values...
func (v *tableView) createRow(row int, callback func(), values ...tableCell) {
	for i, value := range values {
		v.Table.SetCell(row, i, tview.NewTableCell(value.name).SetTextColor(value.color))
	}

	if callback != nil {
		v.selectedFn[row] = callback
	}
}
