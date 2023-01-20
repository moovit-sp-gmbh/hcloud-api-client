package view

import (
	"hcloud-api-client/internal/ui"
	"hcloud-api-client/internal/views"
	"time"

	"github.com/derailed/tview"
)

const (
	splashDelay      = 1 * time.Second
	clusterInfoWidth = 50
	clusterInfoPad   = 15
)

// ExitStatus indicates UI exit conditions.
var ExitStatus = ""

// App represents an application view.
type App struct {
	*ui.App
}

// NewApp returns a K9s app instance.
func NewApp(debug bool) *App {
	a := App{
		App: ui.NewApp(debug),
	}

	return &a
}

// Init initializes the application.
func (a *App) Init() error {
	a.App.Init()
	a.layout()

	return nil
}

// Run starts the application loop.
func (a *App) Run() error {
	go func() {
		<-time.After(splashDelay)
		a.QueueUpdateDraw(func() {
			a.Main.SwitchToPage("main")
			a.LoadContent(views.HCLOUD)
		})
	}()

	if err := a.Application.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) layout() {
	main := tview.NewFlex().SetDirection(tview.FlexRow)

	a.Main.AddPage("main", main, true, false)
	a.Main.AddPage("splash", ui.NewSplash(), true, true)

	main.AddItemAtIndex(0, a.buildHeader(), 7, 1, false)
	main.AddItemAtIndex(1, a.buildCommand(), 3, 1, false)
	main.AddItemAtIndex(2, a.buildContent(), 0, 1, false)
	main.AddItemAtIndex(3, a.buildFooter(), 3, 1, false)

	NewHcloud(a.App)
	NewService(a.App)
	NewIDP(a.App)
	NewHigh5(a.App)
	NewConfig(a.App)
	NewHelp(a.App)
	NewDebug(a.App)
}

func (a *App) buildHeader() tview.Primitive {
	f := tview.NewFlex()
	f.SetDirection(tview.FlexColumn)

	f.AddItem(a.Context(), 0, 1, false)
	f.AddItem(a.Menu(), 0, 1, false)
	f.AddItem(a.Logo(), 41, 2, false)

	return f
}

func (a *App) buildCommand() tview.Primitive {
	f := tview.NewFlex()
	f.SetDirection(tview.FlexColumn)

	f.AddItem(a.Command(), 0, 2, false)

	return f
}

func (a *App) buildContent() tview.Primitive {
	f := tview.NewFlex()
	f.SetDirection(tview.FlexColumn)

	f.AddItem(a.Content(), 0, 1, false)

	return f
}

func (a *App) buildFooter() tview.Primitive {
	f := tview.NewFlex()
	f.SetDirection(tview.FlexColumn)

	f.AddItem(a.Crumbs(), 0, 1, false)

	return f
}
