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
