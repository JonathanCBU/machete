package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Boxed is any tview primitive that also exposes SetInputCapture. Every
// built-in tview widget satisfies this (they all embed *tview.Box), but
// SetInputCapture isn't part of the tview.Primitive interface itself, so
// StaticScreen needs this narrower type to be able to attach escToMenu.
type Boxed interface {
	tview.Primitive
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Box
}

// StaticScreen adapts any single primitive (TextView, Form, Table...) into
// a Screen. Use this for pages that aren't a dynamic list of cards, e.g.
// Home or Logs.
type StaticScreen struct {
	name   string
	hotkey rune
	prim   Boxed
	app    *App
}

func NewStaticScreen(app *App, name string, hotkey rune, prim Boxed) *StaticScreen {
	// escToMenu is the only capture Static screens need — they have no
	// list to navigate, so vimNavCapture wouldn't do anything useful here.
	prim.SetInputCapture(escToMenu(app))
	return &StaticScreen{name: name, hotkey: hotkey, prim: prim, app: app}
}

func (s *StaticScreen) Primitive() tview.Primitive { return s.prim }
func (s *StaticScreen) Name() string               { return s.name }
func (s *StaticScreen) Hotkey() rune               { return s.hotkey }
func (s *StaticScreen) OnShow()                    { s.app.Application().SetFocus(s.prim) }
