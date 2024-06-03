package terminal

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type Terminal struct {
	Screen tcell.Screen
	Width, Height int
	CursorX, CursorY int
	OffsetX, OffsetY int
	logger *slog.Logger
	Footer bool
	History *HistoryItem
	nonVisibleRows int
}

type HistoryItem struct{
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

func (t *Terminal) FillTerminal() {
	t.Screen.Clear()
	defer t.Screen.Show()

	lastRow := t.Height - 1
	if (t.Footer) {
		lastRow -= 1
	}
	for y := 0; y <= lastRow; y++ {
			t.PrintText(0, y, strconv.Itoa(y+t.OffsetY), tcell.StyleDefault)
	}

	if (t.Footer) {
		t.PrintFooter(fmt.Sprintf("Footer. Terminal size: width - %d, height - %d, offset - %d", t.Width, t.Height, t.OffsetY), tcell.StyleDefault)
	}
}

type View interface {

	// Method to Print your View on terminal layout
	Print(*Terminal) int

	KeyUp()
	KeyDown()

	// // Select next element in your View by keyDown
	// SelectNext()

	// // Select previus element in your View by keyUp
	// SelectPrev()

	// Action on selected element by keyEnter
	DoSelected(*Terminal)
}

func (t *Terminal) CurrentView() View {
	return t.History.Item
}

func (t *Terminal) PrintCurrentView() {
	t.Screen.Clear()
	defer t.Screen.Show()

	t.nonVisibleRows = t.CurrentView().Print(t)
}

func (t *Terminal) PrintText(x,y int, text string, style tcell.Style) {
	for _, c := range text {
		t.Screen.SetContent(x, y, c, nil, style)
		x++
	}
}

func (t *Terminal) PrintFooter(footer string, style tcell.Style) {
	t.PrintText(0, t.Height - 1, footer, style)
}

