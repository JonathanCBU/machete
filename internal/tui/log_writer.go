package tui

import (
	"github.com/rivo/tview"
)

// LogWriter is an io.Writer you can pass to log.SetOutput so standard log
// calls show up inside the TUI's Logs screen instead of disappearing into
// the alternate screen buffer. Safe to write to from any goroutine —
// crucially, including before app.Run() has started.
//
// Write() never blocks: it hands the bytes off to a buffered channel and
// returns immediately. A background goroutine drains that channel and
// calls QueueUpdateDraw, which itself blocks until the event loop is
// running — but since that happens on its own goroutine, it can wait
// safely without deadlocking whatever called log.Print.
type LogWriter struct {
	app   *tview.Application
	view  *tview.TextView
	lines chan []byte
}

func NewLogWriter(app *tview.Application, view *tview.TextView) *LogWriter {
	view.SetDynamicColors(true)
	w := &LogWriter{
		app:   app,
		view:  view,
		lines: make(chan []byte, 256), // buffered so early bursts don't drop
	}
	go w.run()
	return w
}

func (w *LogWriter) run() {
	for p := range w.lines {
		p := p
		w.app.QueueUpdateDraw(func() {
			w.view.Write(p)
		})
	}
}

func (w *LogWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	select {
	case w.lines <- buf:
	default:
		// Channel full — drop rather than block the caller. Losing a log
		// line under heavy load beats freezing the app that's trying to
		// log about it.
	}
	return len(p), nil
}
