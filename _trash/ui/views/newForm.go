package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newGenericFunc[age int64 | float64](myAge age) {
	fmt.Println(myAge)
}

func (m *Main) newForm(header string, width int, height int) (*tview.Form, *tview.Grid) {
	form := tview.NewForm()
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetFieldBackgroundColor(tcell.ColorBlue)
	form.SetBorder(true).SetTitle(fmt.Sprintf("   %s   ", header))

	g := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(form, 1, 1, 1, 1, 0, 0, true)

	m.Main.Content.Clear().AddItem(g, 0, 1, false)

	return form, g
}

func (m *Main) newInputField(label string, value *string) *tview.InputField {
	field := tview.NewInputField()
	field.SetLabel(label)
	field.SetChangedFunc(func(text string) {
		*value = text
	})
	return field
}

func (m *Main) newSelectField(label string, options []string, value *string) *tview.DropDown {
	field := tview.NewDropDown()
	field.SetLabel(label)
	field.SetOptions(options, func(text string, index int) {
		*value = text
	})
	return field
}

func (m *Main) newTextArea(label string, width int, height int, value *string) *tview.TextArea {
	field := tview.NewTextArea()
	field.SetLabel(label)
	field.SetText(*value, false)
	field.SetSize(height, width)
	field.SetChangedFunc(func() {
		*value = field.GetText()
	})

	return field
}
