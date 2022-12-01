package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/rivo/tview"
)

type Main struct {
	Root *tview.Flex

	Info *tview.Flex
	Menu *tview.Flex
	Logo *tview.Flex

	Command *tview.Flex

	Content *tview.Flex

	SelectedRow map[int]int

	ActiveUser         *hcloud.User
	ActiveOrganization *hcloud.Organization
}

func (m *Main) Init() *Main {
	m.Root = tview.NewFlex()
	m.Info = tview.NewFlex().AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(false), 0, 1, false)
	m.Menu = tview.NewFlex().AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(false), 0, 1, false)
	m.Logo = tview.NewFlex().AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(false), 0, 1, false)
	m.Command = tview.NewFlex().AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(false), 0, 1, false)
	m.Content = tview.NewFlex().AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(false), 0, 1, false)
	m.SelectedRow = make(map[int]int)

	m.Root.
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(m.Info, 0, 2, false).
				AddItem(m.Menu, 0, 4, false).
				AddItem(m.Logo, 0, 2, false), 0, 1, false).AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(m.Command, 0, 1, false).AddItem(m.Content, 0, 8, false), 0, 5, false), 0, 5, false)

	return m
}
