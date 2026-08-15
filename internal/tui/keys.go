package tui

import "github.com/gdamore/tcell/v2"

// vimNavCapture translates vim-style navigation keys into the standard
// key events tview.List already handles (Up/Down/Home/End), rather than
// reimplementing selection movement ourselves. Attach with
// list.SetInputCapture(vimNavCapture).
//
//	j → Down    k → Up    g → Home (top)    G → End (bottom)
//
// Any other key passes through untouched.
func vimNavCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	case 'g':
		return tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModNone)
	case 'G':
		return tcell.NewEventKey(tcell.KeyEnd, 0, tcell.ModNone)
	default:
		return event
	}
}

// chainCapture runs each capture function in order, short-circuiting as
// soon as one swallows the event (returns nil). Lets a primitive combine
// e.g. vim nav translation with an Esc-to-menu handler without either one
// needing to know about the other.
func chainCapture(fns ...func(*tcell.EventKey) *tcell.EventKey) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		for _, fn := range fns {
			if fn == nil {
				continue
			}
			event = fn(event)
			if event == nil {
				return nil
			}
		}
		return event
	}
}

// escToMenu returns focus to the App's sidebar menu on Esc. Attach it to
// a screen's own primitive (via SetInputCapture, typically chained with
// vimNavCapture) so it only fires while that primitive itself is
// focused — a modal or other overlay that owns Esc for its own purposes
// (closing itself) keeps working, since focus will be on the modal, not
// the underlying screen, when Esc is pressed.
func escToMenu(app *App) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.FocusMenu()
			return nil
		}
		return event
	}
}
