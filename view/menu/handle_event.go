package menu

import (
	"terminal-view/terminal"

	"github.com/gdamore/tcell/v2"
)

func (m *Menu) HandleEvent(ev tcell.Event) terminal.Event {

	switch ev := ev.(type) {
	case *tcell.EventResize:
		m.Resize(ev.Size())
		return terminal.NewResizeEvent(ev.Size())
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			return terminal.NewQuitEvent()
		case tcell.KeyESC:
			return terminal.NewHistoryBackEvent()
		case tcell.KeyUp:
			m.KeyUp()
		case tcell.KeyDown:
			m.KeyDown()
		case tcell.KeyEnter:
			return terminal.NewNextViewEvent(nil)
		}
	case *tcell.EventMouse:
		switch ev.Buttons() {
		case tcell.WheelDown:
			m.WheelDown1(ev)
		case tcell.WheelUp:
			m.WheelUp1(ev)
		}
	case *tcell.EventError:
		return terminal.NewErrorEvent(ev)
	}

	return nil
}
