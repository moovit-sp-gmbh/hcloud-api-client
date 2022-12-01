package ui

import (
	"fmt"
	"hcloud-api-client/config"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
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
func (a *App) getMainTable(title string) *Table {
	// reset menu
	a.resetMenu()

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
		a.showCommand()
	})

	// register global key '?' to open help
	t.registerKey(63, func() {
		a.showHelp()
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

func getMenuTable() *tview.Table {
	table := tview.NewTable()

	table.SetBorder(false)
	table.SetSelectable(false, false)

	return table
}

func addPaddingToString(in string) string {
	return fmt.Sprintf("   %s   ", in)
}

func (a *App) SetFocus(t tview.Primitive) {
	a.Application.SetFocus(t)
}

func (a *App) showConfirm(header string, width int, height int, callback func(bool), text ...string) {
	old := a.Pages.Main.Content.GetItem(0)
	form := tview.NewForm()
	g := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	for _, t := range text {
		form.AddTextView(t, t, 1, 1, false, false)
	}
	form.AddButton("Yes", func() {
		a.Pages.Main.Content.Clear().AddItem(old, 0, 1, false)
		a.SetFocus(old)
		callback(true)
	})
	form.AddButton("No", func() {
		a.Pages.Main.Content.Clear().AddItem(old, 0, 1, false)
		a.SetFocus(old)

		callback(false)
	})
	form.SetBorder(true).SetTitle(addPaddingToString(header))

	a.Pages.Main.Content.Clear().AddItem(g, 0, 1, true)
	a.SetFocus(form)
}

func (a *App) removeFromTable(table *tview.Table, row int, offset int) {
	table.RemoveRow(row)
	table.Select(row-1, 0)
}

func (a *App) showError(err *hcloud.ErrorResponse) {
	textView := tview.NewTextView().SetText(fmt.Sprintf("   %s", err.ToString())).SetTextColor(tcell.ColorRed)
	a.Pages.Main.Command.Clear().AddItem(textView, 0, 1, false)

	time.AfterFunc(5*time.Second, func() {
		textView.SetText("")
	})
}

func (a *App) showCommand() {
	if config.Config.Active != nil {
		textInput := tview.NewInputField().SetLabel(":").SetAutocompleteFunc(func(currentText string) (entries []string) {
			return []string{"home", "ctx", "context", "svc", "service", "idp", "high5", "quit", "help"}
		})
		textInput.SetAutocompletedFunc(func(text string, index, source int) bool {
			a.autoCompleteDone(text)
			a.Pages.Main.Command.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
			return true
		})

		textInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyUp {
				return nil
			}
			if event.Key() == tcell.KeyDown {
				return nil
			}
			return event
		})

		textInput.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				a.autoCompleteDone(textInput.GetText())
				a.Pages.Main.Command.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
			}
		})
		a.Pages.Main.Command.Clear().AddItem(textInput, 0, 1, false)
		a.SetFocus(textInput)
	}
}

func (a *App) autoCompleteDone(text string) {
	switch text {
	case "home":
		a.buildHcloud()
	case "svc", "service":
		a.buildServices()
	case "ctx", "context":
		a.buildConfigContext(false)
	case "idp":
		a.buildIdp()
	case "high5":
		a.buildHigh5Apps()
	case "help":
		a.showHelp()
	case "quit":
		a.Application.Stop()
	}
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
func (a *App) initLastSelection(t *Table, uniqueClassId int, defaultSelection int) {
	defaultSelection = int(math.Min(float64(defaultSelection), float64(t.GetRowCount()-1)))
	t.Select(defaultSelection, 0)

	if uniqueClassId > -1 {
		// save last selection
		t.SetSelectionChangedFunc(func(row, column int) {
			a.Pages.Main.SelectedRow[uniqueClassId] = row
		})

		// load last selection
		if row, ok := a.Pages.Main.SelectedRow[uniqueClassId]; ok {
			if row < t.GetRowCount() {
				t.Select(row, 0)
			}
		}
	}
}

func (a *App) resetMenu() {
	a.Pages.Main.Menu.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
}
