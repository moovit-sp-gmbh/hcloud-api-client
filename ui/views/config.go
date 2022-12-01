package views

import (
	"hcloud-api-client/pkg"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/skratchdot/open-golang/open"
)

func (m *Main) buildConfig() {
	table := m.getMainTable("helmut.cloud > config")

	table.SetCell(0, 0, tview.NewTableCell("..").SetTextColor(tcell.ColorWhite))
	table.SetCell(1, 0, tview.NewTableCell("context").SetTextColor(tcell.ColorWhite))
	table.SetCell(2, 0, tview.NewTableCell("open config folder").SetTextColor(tcell.ColorWhite))
	table.SetCell(3, 0, tview.NewTableCell("open config file").SetTextColor(tcell.ColorWhite))

	m.initLastSelection(table, main_config, 1)

	table.SetSelectedFunc(func(row int, column int) {
		switch row {
		case 0:
			m.BuildHcloud()
		case 1:
			m.buildContext(false)
		case 2:
			open.Run(pkg.GetHomeDir() + "/.hcloud/")
		case 3:
			open.Run(pkg.GetHomeDir() + "/.hcloud/config.yml")
		}
	})

	table.registerKey(27, func() {
		m.BuildHcloud()
	})

	m.Main.Content.Clear().AddItem(table, 0, 1, false)
	m.Application.SetFocus(table)
}
