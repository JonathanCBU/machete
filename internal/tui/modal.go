package tui

import "github.com/rivo/tview"

// centerModal wraps any primitive in nested Flex spacers so it renders as
// a fixed-size box centered on screen, rather than filling it. This is the
// standard tview technique for building custom modals out of any
// primitive (not just tview.Modal's message+buttons).
func centerModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
