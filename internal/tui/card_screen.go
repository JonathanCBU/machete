package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CardRenderer turns a domain item of type T into display text: the list
// row's title/summary, and the longer text shown in the modal when the
// row is opened. This is the one piece of domain knowledge CardScreen
// needs — everything else about layout, navigation, and the modal is
// generic and reusable across item types.
type CardRenderer[T any] func(item T) (title, summary, detail string)

// CardScreen is a navigable list of "cards" that expand into a modal with
// more detail on Enter. Instantiate it with any T (profiles, servers,
// alerts...) plus a CardRenderer[T]. Because it's generic, the tui package
// never needs to know what a "Profile" is.
type CardScreen[T any] struct {
	name   string
	hotkey rune
	render CardRenderer[T]
	list   *tview.List
	app    *tview.Application
	pages  *tview.Pages
}

func NewCardScreen[T any](
	app *tview.Application,
	pages *tview.Pages,
	name string,
	hotkey rune,
	items []T,
	render CardRenderer[T],
) *CardScreen[T] {
	s := &CardScreen[T]{
		name:   name,
		hotkey: hotkey,
		render: render,
		app:    app,
		pages:  pages,
		list:   tview.NewList().ShowSecondaryText(true),
	}
	s.list.SetBorder(true).SetTitle(" " + name + " ")
	s.list.SetInputCapture(vimNavCapture)
	s.SetItems(items)
	return s
}

// SetItems (re)populates the list. Call this any time — after an API
// call, a timer tick, a filter change — to dynamically regenerate the
// screen's contents without rebuilding the screen itself.
func (s *CardScreen[T]) SetItems(items []T) {
	s.list.Clear()
	for _, item := range items {
		title, summary, detail := s.render(item)
		s.list.AddItem(title, summary, 0, func() {
			s.showModal(title, detail)
		})
	}
}

func (s *CardScreen[T]) showModal(title, detail string) {
	text := tview.NewTextView().SetText(detail)
	text.SetBorder(true).SetTitle(" " + title + " ")

	modalName := s.name + "-modal"
	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			s.pages.RemovePage(modalName)
			s.app.SetFocus(s.list)
			return nil
		}
		return event
	})

	s.pages.AddPage(modalName, centerModal(text, 40, 10), true, true)
	s.app.SetFocus(text)
}

func (s *CardScreen[T]) Primitive() tview.Primitive { return s.list }
func (s *CardScreen[T]) Name() string               { return s.name }
func (s *CardScreen[T]) Hotkey() rune               { return s.hotkey }
func (s *CardScreen[T]) OnShow()                    { s.app.SetFocus(s.list) }
