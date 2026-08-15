package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App wires together the sidebar menu, the Pages content area, and the
// hotkey routing between them. Screens are registered dynamically, so the
// menu and key handling never need hardcoded switch statements — adding
// a screen in main.go is enough to make it navigable.
type App struct {
	app     *tview.Application
	pages   *tview.Pages
	menu    *tview.List
	root    *tview.Flex
	screens []Screen
	byKey   map[rune]Screen
}

func NewApp() *App {
	a := &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
		menu:  tview.NewList().ShowSecondaryText(false),
		byKey: map[rune]Screen{},
	}
	a.menu.SetBorder(true).SetTitle(" Menu (q: quit) ")
	// The menu is navigable the same way card screens are: j/k/g/G plus
	// arrows via List's own defaults. Esc has no meaning on the menu
	// itself, so it isn't chained here (unlike on the content screens).
	a.menu.SetInputCapture(vimNavCapture)

	a.root = tview.NewFlex().
		AddItem(a.menu, 20, 0, false).
		AddItem(a.pages, 0, 1, true)
	return a
}

// Application and Pages expose the underlying primitives so main.go can
// construct screens (which need references to focus/switch pages) without
// tui needing to know about those screens' concrete types.
func (a *App) Application() *tview.Application { return a.app }
func (a *App) Pages() *tview.Pages             { return a.pages }

// FocusMenu moves focus to the sidebar. Screens wire this to Esc via
// escToMenu so pressing Esc from a screen's own primitive returns here.
func (a *App) FocusMenu() {
	a.app.SetFocus(a.menu)
}

// switchTo makes s the visible page and focuses its content. Shared by
// both the menu's Enter handler and the global hotkey capture so there's
// one place that defines what "switching to a screen" means.
func (a *App) switchTo(s Screen) {
	a.pages.SwitchToPage(s.Name())
	s.OnShow()
}

// Register adds a screen: wires its hotkey, adds it to Pages, and adds a
// corresponding row to the menu List. The first registered screen is
// shown by default.
func (a *App) Register(s Screen) {
	a.screens = append(a.screens, s)
	a.byKey[s.Hotkey()] = s
	a.pages.AddPage(s.Name(), s.Primitive(), true, len(a.screens) == 1)

	// The shortcut rune makes List show e.g. "h) Home" and also lets
	// pressing that rune while the menu is focused jump straight to it —
	// on top of the global hotkey capture in Run, which works from
	// anywhere regardless of focus.
	a.menu.AddItem(s.Name(), "", s.Hotkey(), func() {
		a.switchTo(s)
	})
}

func (a *App) Run() error {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			a.app.Stop()
			return nil
		}
		if s, ok := a.byKey[event.Rune()]; ok {
			a.switchTo(s)
			return nil
		}
		return event
	})
	// Start focused on the menu — it's now the natural entry point:
	// navigate with j/k, Enter into a screen, Esc back out to here.
	return a.app.SetRoot(a.root, true).SetFocus(a.menu).Run()
}
