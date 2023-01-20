package views

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	*tview.Table

	keys map[int32]func()
}

func (t *Table) registerKey(key int32, fn func()) *Table {
	t.keys[key] = fn
	return t
}

// getMainTable returns a default table
func (m *Main) getMainTable(title string) *Table {
	// reset menu
	m.resetMenu()

	table := tview.NewTable()

	table.SetBorder(true)
	table.SetBorderPadding(1, 1, 1, 1)
	table.SetTitle(fmt.Sprintf("   %s   ", title))
	table.SetTitleAlign(0)
	table.SetSelectable(true, false)
	table.SetBorderColor(tcell.ColorBlue)

	keys := make(map[int32]func())

	t := &Table{table, keys}

	// register global key ':' to open command prompt
	t.registerKey(58, func() {
		m.buildCommand()
	})

	// register global key '?' to open help
	t.registerKey(63, func() {
		m.buildHelp()
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if fn, ok := keys[event.Rune()]; ok {
			fn()
		} else if fn, ok := keys[int32(event.Key())]; ok {
			fn()
		}
		return event
	})

	return t
}

const (
	main_hcloud = iota

	main_help

	main_config
	main_config_context

	main_services

	main_idp
	main_idp_organizations

	main_high5
	main_high5_apps
	main_high5_app_detail
	main_high5_app_events
	main_high5_app_event_streams
	main_high5_app_settings
	main_high5_app_settings_webhooks
	main_high5_app_settings_webhook_logs
)

// initLastSelection loads the default selection
// if uniqueClassId is greater than -1, will save future selections to memory and load last value from memory
func (m *Main) initLastSelection(t *Table, uniqueClassId int, defaultSelection int) {
	defaultSelection = int(math.Min(float64(defaultSelection), float64(t.GetRowCount()-1)))
	t.Select(defaultSelection, 0)

	if uniqueClassId > -1 {
		// save last selection
		t.SetSelectionChangedFunc(func(row, column int) {
			m.Main.SelectedRow[uniqueClassId] = row
		})

		// load last selection
		if row, ok := m.Main.SelectedRow[uniqueClassId]; ok {
			if row < t.GetRowCount() {
				t.Select(row, 0)
			}
		}
	}
}
