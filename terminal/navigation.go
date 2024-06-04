package terminal

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

func (t *Terminal) KeyUp() {
	t.CurrentView().KeyUp()
}

func (t *Terminal) KeyDown() {
	t.CurrentView().KeyDown()
}

func (t *Terminal) WheelUp(ev *tcell.EventMouse) {
	x, y := ev.Position()
	t.CurrentView().WheelUp(&MouseEvent{x: x, y: y})
}

func (t *Terminal) WheelDown(ev *tcell.EventMouse) {
	x, y := ev.Position()
	t.CurrentView().WheelDown(&MouseEvent{x: x, y: y})
}

// Action on KeyEnter
func (t *Terminal) DoSelected() {
	t.CurrentView().DoSelected(t)
}

var (
	ErrBrokenHistory = errors.New("no view history")
	ErrLastItem      = errors.New("it's a last item in history")
)

// Action on KeyESC
//
// If CurrentView it is a last View in History return ErrLastItem
func (t *Terminal) GoToPrevView() error {
	if t.History == nil {
		return ErrBrokenHistory
	}

	if t.History.Prev == nil {
		return ErrLastItem
	}

	t.History = t.History.Prev
	return nil

}

func (t *Terminal) GoToNewView(view View) {
	item := createHistoryItem(view)

	if t.History == nil {
		t.History = item
	} else {
		item.Prev = t.History
		t.History = item
	}

}

func createHistoryItem(view View) *HistoryItem {
	return &HistoryItem{
		Item: view,
	}
}
