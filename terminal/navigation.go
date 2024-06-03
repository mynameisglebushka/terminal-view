package terminal

import "errors"

func (t *Terminal) ScrollUp() {
	if t.OffsetY > 0 {
		t.OffsetY -= 1
	}
}

func (t *Terminal) ScrollDown() {
	if t.nonVisibleRows > 0 {
		t.OffsetY += 1
	}
}

// Action on KeyUp
// func (t *Terminal) SelectPrev() {
// 	t.CurrentView().SelectPrev()
// }

func (t *Terminal) KeyUp() {
	t.CurrentView().KeyUp()
}

// Action on KeyDown
// func (t *Terminal) SelectNext() {
// 	t.CurrentView().SelectNext()
// }

func (t *Terminal) KeyDown() {
	t.CurrentView().KeyDown()
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
