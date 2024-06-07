package terminal

import (
	"errors"
)

var (
	ErrBrokenHistory = errors.New("no view history")
	ErrLastItem      = errors.New("it's a last item in history")
)

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

func (t *Terminal) GoToNextView(view View) {
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
