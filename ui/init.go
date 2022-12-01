package ui

import (
	"hcloud-api-client/ui/pages"
	"hcloud-api-client/ui/views"

	"github.com/moovit-sp-gmbh/hcloud-sdk-go"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application

	Pages Pages
}

type Pages struct {
	Main *pages.Main
}

func NewApp() *App {
	a := App{
		Application: tview.NewApplication(),
		Pages:       Pages{Main: &pages.Main{}},
	}

	a.Pages.Main.Init()

	return &a
}

var activeUser *hcloud.User
var activeOrganization *hcloud.Organization

func (a *App) Init() {
	main := views.Main{Application: a.Application, Main: a.Pages.Main}
	main.BuildLogo()
	main.BuildHcloud()

	if err := a.Application.SetRoot(a.Pages.Main.Root, true).SetFocus(a.Pages.Main.Content.GetItem(0)).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
