package tui

import "github.com/rivo/tview"

// StaticScreen adapts any single primitive (TextView, Form, Table...) into
// a Screen. Use this for pages that aren't a dynamic list of cards, e.g.
// Home or Logs.
type StaticScreen struct {
	name   string
	hotkey rune
	prim   tview.Primitive
	app    *tview.Application
}

func NewStaticScreen(app *tview.Application, name string, hotkey rune, prim tview.Primitive) *StaticScreen {
	return &StaticScreen{name: name, hotkey: hotkey, prim: prim, app: app}
}

func (s *StaticScreen) Primitive() tview.Primitive { return s.prim }
func (s *StaticScreen) Name() string               { return s.name }
func (s *StaticScreen) Hotkey() rune               { return s.hotkey }
func (s *StaticScreen) OnShow()                    { s.app.SetFocus(s.prim) }
