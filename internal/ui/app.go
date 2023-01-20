package ui

import (
	"fmt"
	"hcloud-api-client/internal/views"
	"strings"
	"time"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	*tview.Application

	Main       *Pages
	keyActions KeyActions
	views      map[string]tview.Primitive

	debug     bool
	debugLogs []DebugLog
}

type DebugLog struct {
	Timestamp int64
	Message   string
}

// NewApp returns a new app.
func NewApp(debug bool) *App {
	a := App{
		Application: tview.NewApplication(),
		keyActions:  make(KeyActions),
		Main:        NewPages(),

		debug:     debug,
		debugLogs: []DebugLog{},
	}

	a.views = map[string]tview.Primitive{
		"context": NewContext(),
		"menu":    NewMenu(),
		"logo":    NewLogo(),
		"command": NewCommand(),
		"content": NewContent(),
		"crumbs":  NewCrumbs(),
		"prompt":  NewPrompt(&a),
	}

	return &a
}

// LoadContent loads an existing page within the Content page.
func (a *App) LoadContent(viewId int) *App {
	a.Content().SwitchToPage(views.Views[viewId].Name)

	return a
}

// Init initializes the application.
func (a *App) Init() {
	a.bindKeys()

	a.SetRoot(a.Main, true).EnableMouse(true)
}

// Context return the app context.
func (a *App) Context() *Context {
	return a.views["context"].(*Context)
}

// Menu returns the app menu.
func (a *App) Menu() *Menu {
	return a.views["menu"].(*Menu)
}

// Logo return the app logo.
func (a *App) Logo() *Logo {
	return a.views["logo"].(*Logo)
}

// Command return the app context.
func (a *App) Command() *Command {
	return a.views["command"].(*Command)
}

// Content return the app context.
func (a *App) Content() *Content {
	return a.views["content"].(*Content)
}

// Crumbs returns the app Crumbs.
func (a *App) Crumbs() *Crumbs {
	return a.views["crumbs"].(*Crumbs)
}

// Prompt returns the app prompt.
func (a *App) Prompt() *Prompt {
	return a.views["prompt"].(*Prompt)
}

func (a *App) bindKeys() {
	// a.keyActions[tcell.KeyCtrlK] = NewKeyAction("ctrl+k", "command", a.activateCmd)
	// a.keyActions[tcell.KeyCtrlC] = NewKeyAction("ctrl+c", "Quit", a.quitCmd)
}

func (a *App) activateCmd(evt *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (a *App) quitCmd(evt *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (a *App) IsDebugEnabled() bool {
	return a.debug
}

func (a *App) Debug(pattern string, values ...interface{}) *App {
	pattern = strings.ReplaceAll(pattern, "%s", fmt.Sprintf("[#%06x::b]%s[#%06x::b]", tcell.ColorDarkOrange.TrueColor().Hex(), "%s", tcell.ColorWhite.TrueColor().Hex()))
	pattern = strings.ReplaceAll(pattern, "%d", fmt.Sprintf("[#%06x::b]%s[#%06x::b]", tcell.ColorBlue.TrueColor().Hex(), "%d", tcell.ColorWhite.TrueColor().Hex()))

	a.debugLogs = append(a.debugLogs, DebugLog{Timestamp: time.Now().UnixNano(), Message: fmt.Sprintf(pattern, values...)})
	return a
}

func (a *App) GetDebugLogs() []DebugLog {
	return a.debugLogs
}
