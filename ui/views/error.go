package views

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/rivo/tview"
)

func (m *Main) showError(err *hcloud.ErrorResponse) {
	textView := tview.NewTextView().SetText(fmt.Sprintf("   %s", err.ToString())).SetTextColor(tcell.ColorRed)
	m.Main.Command.Clear().AddItem(textView, 0, 1, false)

	time.AfterFunc(5*time.Second, func() {
		textView.SetText("")
	})
}
