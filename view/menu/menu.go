package menu

import (
	"fmt"
	"strings"
	"terminal-view/terminal"

	"github.com/gdamore/tcell/v2"
)

type Menu struct {
	Title        string
	FirstItem    *MenuItem
	LastItem     *MenuItem
	SelectedItem *MenuItem
}

func NewMenu(title string) *Menu {
	return &Menu{
		Title: title,
	}
}

func (m *Menu) AddItem(item *MenuItem) *Menu {
	if m.FirstItem != nil {
		item.Parent = m.LastItem
		item.Child = m.FirstItem
		m.LastItem.Child = item
		m.FirstItem.Parent = item
		m.LastItem = item
	} else {
		m.FirstItem = item
		m.LastItem = item
		m.SelectedItem = item
		m.FirstItem.IsSelected = true
	}
	return m
}

func (m *Menu) KeyUp() {
	m.selectParent()
}

func (m *Menu) selectParent() {
	if m.SelectedItem.Parent != nil {
		m.SelectedItem.Parent.IsSelected = true
		m.SelectedItem.IsSelected = false
		m.SelectedItem = m.SelectedItem.Parent
	}
}

func (m *Menu) KeyDown() {
	m.selectChild()
}

func (m *Menu) selectChild() {
	if m.SelectedItem.Child != nil {
		m.SelectedItem.Child.IsSelected = true
		m.SelectedItem.IsSelected = false
		m.SelectedItem = m.SelectedItem.Child
	}
}

func (m *Menu) Print(t *terminal.Terminal) int {
	var (
		lastRow     int       = t.Height - 1
		item        *MenuItem = m.FirstItem
		countOffset int       = 0
	)

	if t.Footer {
		lastRow -= 1
	}

	// Куда ведет клавиша ESC
	t.PrintText(0, 0, fmt.Sprintf("Меню: %s. Для выхода нажмите ESC", m.Title), tcell.StyleDefault)

	// Список элементов 1/3 - Описание выбранного элемента 2/3
	// Находим пункт меню в соответствии с отступом
	for countOffset < t.OffsetY {
		countOffset++
		item = item.Child
	}

	for y := 2; y < lastRow; y++ {
		t.PrintText(0, y, item.Title, item.UnderlineIfSelected())
		if item.Child == nil || item.Child == m.FirstItem {
			break
		}
		item = item.Child
	}

	for y := 2; y < lastRow; y++ {
		t.PrintText(t.Width/3, y, " ", tcell.StyleDefault.Background(tcell.ColorGrey))
	}

	spltDesc := strings.Split(m.SelectedItem.Description, "\n")
	i := 0
	for y := 2; y < lastRow; y++ {
		if i >= len(spltDesc) {
			break
		}
		t.PrintText(t.Width/3+2, y, spltDesc[i], tcell.StyleDefault)
		i++
	}

	var countNoPrintedRows int
	for item.Child != m.FirstItem {
		countNoPrintedRows++
		item = item.Child
	}

	// Футер - выбранный элемент
	if t.Footer {
		t.PrintFooter(fmt.Sprintf("Selected: %s", m.SelectedItem.Title), tcell.StyleDefault)
	}

	return countNoPrintedRows
}

func (m *Menu) DoSelected(t *terminal.Terminal) {
	if m.SelectedItem != nil {
		if m.SelectedItem.SubMenu != nil {
			t.GoToNewView(m.SelectedItem.SubMenu)
		}
	}
}

type MenuItem struct {
	Title       string
	Description string
	Parent      *MenuItem
	Child       *MenuItem
	SubMenu     *Menu
	IsSelected  bool
}

func NewMenuItem(title string) *MenuItem {
	return &MenuItem{
		Title: title,
	}
}

func (mi *MenuItem) WithSubMenu(menu *Menu) *MenuItem {
	mi.SubMenu = menu
	return mi
}

func (mi *MenuItem) UnderlineIfSelected() tcell.Style {
	if mi.IsSelected {
		return tcell.StyleDefault.Underline(true)
	}
	return tcell.StyleDefault.Underline(false)
}
