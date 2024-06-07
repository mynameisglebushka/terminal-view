package terminal

import "time"

func (t *Terminal) Execute() error {

eventloop:
	for {
		event := t.CurrentView().HandleEvent(t.Screen.PollEvent())
		switch ev := event.(type) {
		case *ResizeEvent:
			t.Resize()
		case *NextViewEvent:
			if ev.view != nil {
				t.GoToNextView(ev.view)
			}
		case *HistoryBackEvent:
			if t.GoToPrevView() != nil {
				break eventloop
			}
		case *FocusEvent:
			t.Focused = ev.inFocus
		case *QuitEvent:
			break eventloop
		case *ErrorEvent:
			panic(ev)
		}

		t.PrintCurrentView()
	}
	t.logger.Info("terminal session over", "time", time.Since(t.startupTime))
	return nil
}
