package tui

import "github.com/rivo/tview"

// Screen is anything that can be registered with App as a right-side page
// and switched to via a sidebar hotkey. Keeping this minimal means new
// screen kinds (static text, card lists, forms, tables...) only need to
// satisfy four methods to plug into the same menu/routing machinery.
type Screen interface {
	Primitive() tview.Primitive // what Pages should display
	Name() string               // page id + menu label
	Hotkey() rune               // key that switches to this screen
	OnShow()                    // called after becoming visible (e.g. set focus)
}
