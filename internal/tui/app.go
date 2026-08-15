package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App wires together the sidebar menu, the Pages content area, and the
// hotkey routing between them. Screens are registered dynamically, so the
// menu text and key handling never need hardcoded switch statements —
// adding a screen in main.go is enough to make it navigable.
type App struct {
	app     *tview.Application
	pages   *tview.Pages
	menu    *tview.TextView
	root    *tview.Flex
	screens []Screen
	byKey   map[rune]Screen
}

func NewApp() *App {
	a := &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
		menu:  tview.NewTextView().SetDynamicColors(true),
		byKey: map[rune]Screen{},
	}
	a.menu.SetBorder(true).SetTitle(" Menu ")
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

// Register adds a screen, wires its hotkey, and refreshes the menu text.
// The first registered screen is shown by default.
func (a *App) Register(s Screen) {
	a.screens = append(a.screens, s)
	a.byKey[s.Hotkey()] = s
	a.pages.AddPage(s.Name(), s.Primitive(), true, len(a.screens) == 1)
	a.rebuildMenu()
}

func (a *App) rebuildMenu() {
	text := "[yellow]Navigation[-]\n\n"
	for _, s := range a.screens {
		text += fmt.Sprintf("[green]%c[-]  %s\n", s.Hotkey(), s.Name())
	}
	text += "\n[yellow]General[-]\n\n[green]q[-]  Quit\n"
	a.menu.SetText(text)
}

func (a *App) Run() error {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			a.app.Stop()
			return nil
		}
		if s, ok := a.byKey[event.Rune()]; ok {
			a.pages.SwitchToPage(s.Name())
			s.OnShow()
			return nil
		}
		return event
	})
	return a.app.SetRoot(a.root, true).SetFocus(a.pages).Run()
}
