package views

import (
	"hcloud-api-client/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *Main) buildCommand() {
	if config.Config.Active != nil {
		textInput := tview.NewInputField().SetLabel(":").SetAutocompleteFunc(func(currentText string) (entries []string) {
			return []string{"home", "ctx", "context", "svc", "service", "idp", "high5", "quit", "help"}
		})
		textInput.SetAutocompletedFunc(func(text string, index, source int) bool {
			autoCompleteDone(m, text)
			m.Main.Command.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
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
				autoCompleteDone(m, textInput.GetText())
				m.Main.Command.Clear().AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
			}
		})
		m.Main.Command.Clear().AddItem(textInput, 0, 1, false)
		m.SetFocus(textInput)
	}
}

func autoCompleteDone(m *Main, text string) {
	switch text {
	case "home":
		m.BuildHcloud()
	case "svc", "service":
		m.buildServices()
	case "ctx", "context":
		m.buildContext(false)
	case "idp":
		m.buildIdp()
	case "high5":
		m.buildHigh5()
	case "help":
		m.buildHelp()
	case "quit":
		m.Application.Stop()
	}
}
