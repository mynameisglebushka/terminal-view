package terminal

import "github.com/gdamore/tcell/v2"

func (t *Terminal) Execute() error {
	eventloop:
	for {
		ev := t.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			t.Resize()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				break eventloop
			case tcell.KeyESC:
				if t.GoToPrevView() != nil {
					break eventloop
				}
			case tcell.KeyUp:
				// t.SelectPrev()
				t.KeyUp()
			case tcell.KeyDown:
				// t.SelectNext()
				t.KeyDown()
			case tcell.KeyEnter:
				t.DoSelected()

			}
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.WheelDown:
				t.ScrollDown()
			case tcell.WheelUp:
				t.ScrollUp()
			}
		case *tcell.EventError:
			panic(ev)
		}

		t.PrintCurrentView()
	}
	return nil
}