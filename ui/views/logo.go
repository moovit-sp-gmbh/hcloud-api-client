package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var logoSmall = []string{
	`.__           .__                   .___  `,
	`|  |__   ____ |  |   ____  __ __  __| _/  `,
	`|  |  \_/ ___\|  |  /  _ \|  |  \/ __ |   `,
	`|   Y  \  \___|  |_(  <_> )  |  / /_/ |   `,
	`|___|  /\___  >____/\____/|____/\____ |   `,
}

var logoBig = []string{
	`.__           .__                  __           .__                   .___`,
	`|  |__   ____ |  |   _____  __ ___/  |_    ____ |  |   ____  __ __  __| _/`,
	`|  |  \_/ __ \|  |  /     \|  |  \   __\ _/ ___\|  |  /  _ \|  |  \/ __ |`,
	`|   Y  \  ___/|  |_|  Y Y  \  |  /|  |   \  \___|  |_(  <_> )  |  / /_/ |`,
	`|___|  /\___  >____/__|_|  /____/ |__|   /\___  >____/\____/|____/\____ |`,
}

func (m *Main) BuildLogo() {
	m.Main.Logo.Clear().AddItem(tview.NewTextView().SetText(strings.Join(logoSmall, "\n")).SetTextColor(tcell.ColorOrangeRed).SetTextAlign(tview.AlignRight), 0, 1, false)
}
