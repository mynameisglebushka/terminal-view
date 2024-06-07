package terminal

import "time"

type Event interface {
	When() time.Time
}

type ResizeEvent struct {
	when time.Time
	x, y int
}

func NewResizeEvent(x, y int) *ResizeEvent {
	return &ResizeEvent{
		when: time.Now(),
		x:    x,
		y:    y,
	}
}

func (e *ResizeEvent) When() time.Time {
	return e.when
}

func (e *ResizeEvent) Size() (int, int) {
	return e.x, e.y
}

type ErrorEvent struct {
	when  time.Time
	error error
}

func NewErrorEvent(err error) *ErrorEvent {
	return &ErrorEvent{
		when:  time.Now(),
		error: err,
	}
}

func (e *ErrorEvent) When() time.Time {
	return e.when
}

type QuitEvent struct {
	when time.Time
}

func NewQuitEvent() *QuitEvent {
	return &QuitEvent{
		when: time.Now(),
	}
}

func (e *QuitEvent) When() time.Time {
	return e.when
}

type HistoryBackEvent struct {
	when time.Time
}

func NewHistoryBackEvent() *HistoryBackEvent {
	return &HistoryBackEvent{
		when: time.Now(),
	}
}

func (e *HistoryBackEvent) When() time.Time {
	return e.when
}

type NextViewEvent struct {
	when time.Time
	view View
}

func NewNextViewEvent(view View) *NextViewEvent {
	return &NextViewEvent{
		when: time.Now(),
		view: view,
	}
}

func (e *NextViewEvent) When() time.Time {
	return e.when
}

type FocusEvent struct {
	when    time.Time
	inFocus bool
}

func NewFocusEvent(ok bool) *FocusEvent {
	return &FocusEvent{
		when:    time.Now(),
		inFocus: ok,
	}
}

func (e *FocusEvent) When() time.Time {
	return e.when
}
