package terminal

import (
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Terminal struct {
	Screen           tcell.Screen
	Width, Height    int
	logger           *slog.Logger
	Footer           bool
	History          *HistoryItem
	nonVisibleRows   int
}

type HistoryItem struct {
	Prev *HistoryItem
	Item View
}

type TerminalFunc func(*Terminal)

func WithLogger(log *slog.Logger) TerminalFunc {
	return func(t *Terminal) {
		t.logger = log
	}
}

func WithFooter() TerminalFunc {
	return func(t *Terminal) {
		t.Footer = true
	}
}

func defaultTerminal() *Terminal {
	return &Terminal{
		logger: slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		),
	}
}

func NewTerminal(opts ...TerminalFunc) (*Terminal, error) {

	t := defaultTerminal()
	for i := range opts {
		opts[i](t)
	}

	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = s.Init(); err != nil {
		return nil, err
	}

	s.EnableMouse(tcell.MouseMotionEvents)

	w, h := s.Size()

	t.Screen = s
	t.Width = w
	t.Height = h

	return t, nil
}

func (t *Terminal) Recover() {
	p := recover()
	if p != nil {
		t.logger.Error("recover", "err", p)
	}
	t.CloseTerminal()
}

func (t *Terminal) CloseTerminal() {
	t.Screen.Fini()
}

func (t *Terminal) Resize() {
	w, h := t.Screen.Size()
	t.Width = w
	t.Height = h
}

func (t *Terminal) AddStartView(view View) {
	t.GoToNewView(view)
}

func (t *Terminal) CurrentView() View {
	return t.History.Item
}

type MouseEvent struct {
	x, y int
}

func (me *MouseEvent) XY() (int, int) {
	return me.x, me.y
}

type View interface {
	KeyUp()

	KeyDown()

	WheelUp(*MouseEvent)

	WheelDown(*MouseEvent)

	// Method to Print your View on terminal layout
	Print(*Terminal) int

	// Action on selected element by keyEnter
	DoSelected(*Terminal)
}

func (t *Terminal) PrintCurrentView() {
	t.Screen.Clear()
	defer t.Screen.Show()

	t.nonVisibleRows = t.CurrentView().Print(t)
}

func (t *Terminal) PrintText(x, y int, text string, style tcell.Style) {
	for _, c := range text {
		t.Screen.SetContent(x, y, c, nil, style)
		x++
	}
}

func (t *Terminal) PrintRune(x, y int, ch rune, style tcell.Style) {
	t.Screen.SetContent(x, y, ch, nil, style)
}

func (t *Terminal) PrintFooter(footer string, style tcell.Style) {
	t.PrintText(0, t.Height-1, footer, style)
}
