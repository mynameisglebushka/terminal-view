package menu

import (
	"fmt"
	"terminal-view/terminal"

	"github.com/gdamore/tcell/v2"
)

type Menu struct {
	Title        string
	firstItem    *MenuItem
	lastItem     *MenuItem
	selectedItem *MenuItem

	offsetY        int
	nonVisibleRows int
}

func NewMenu(title string) *Menu {
	return &Menu{
		Title: title,
	}
}

func (m *Menu) AddItem(item *MenuItem) *Menu {
	if m.firstItem != nil {
		item.parent = m.lastItem
		item.child = m.firstItem
		m.lastItem.child = item
		m.firstItem.parent = item
		m.lastItem = item
	} else {
		m.firstItem = item
		m.lastItem = item
		m.selectedItem = item
		m.firstItem.isSelected = true
	}
	return m
}

func (m *Menu) KeyUp() {
	m.selectParent()
}

func (m *Menu) selectParent() {
	if m.selectedItem.parent != nil {
		m.selectedItem.parent.isSelected = true
		m.selectedItem.isSelected = false
		m.selectedItem = m.selectedItem.parent
	}
}

func (m *Menu) KeyDown() {
	m.selectChild()
}

func (m *Menu) selectChild() {
	if m.selectedItem.child != nil {
		m.selectedItem.child.isSelected = true
		m.selectedItem.isSelected = false
		m.selectedItem = m.selectedItem.child
	}
}

func (m *Menu) WheelUp(ev *terminal.MouseEvent) {
	if m.offsetY > 0 {
		m.offsetY -= 1
	}

}

func (m *Menu) WheelDown(ev *terminal.MouseEvent) {
	if m.nonVisibleRows > 0 {
		m.offsetY += 1
	}
}

func (m *Menu) Print(t *terminal.Terminal) int {
	var (
		lastRow     int       = t.Height - 1
		item        *MenuItem = m.firstItem
		countOffset int       = 0
	)

	if t.Footer {
		lastRow -= 1
	}

	// Куда ведет клавиша ESC
	t.PrintText(0, 0, fmt.Sprintf("Меню: %s. Для выхода нажмите ESC", m.Title), tcell.StyleDefault)

	// Список элементов 1/3 - Описание выбранного элемента 2/3
	// Находим пункт меню в соответствии с отступом
	for countOffset < m.offsetY {
		countOffset++
		item = item.child
	}

	// Print menu items
	for y := 2; y < lastRow; y++ {
		t.PrintText(0, y, item.Title, item.UnderlineIfSelected())
		if item.child == nil || item.child == m.firstItem {
			break
		}
		item = item.child
	}

	var (
		separatorX = t.Width / 3
	)

	// Print vertical separator
	for y := 2; y < lastRow; y++ {
		t.PrintText(separatorX, y, " ", tcell.StyleDefault.Background(tcell.ColorGrey))
	}

	// Print description of selected item
	// var (
	// 	spltDesc = strings.Split(m.selectedItem.Description, "\n")
	// 	i = 0
	// )
	// for y := 2; y < lastRow; y++ {
	// 	if i >= len(spltDesc) {
	// 		break
	// 	}
	// 	t.PrintText(separatorX+2, y, spltDesc[i], tcell.StyleDefault)
	// 	i++
	// }

	var (
		descriptionXStart     = separatorX + 2
		x                 int = descriptionXStart
		y                 int
	)
	for _, ch := range m.selectedItem.Description {
		if x >= t.Width {
			y++
			x = descriptionXStart
		}
		if ch == 10 {
			y++
			x = descriptionXStart
		}
		t.PrintRune(x, y, ch, tcell.StyleDefault)
		x++
	}

	var countNoPrintedRows int
	for item.child != m.firstItem {
		countNoPrintedRows++
		item = item.child
	}

	// Футер - выбранный элемент
	if t.Footer {
		t.PrintFooter(fmt.Sprintf("Selected: %s", m.selectedItem.Title), tcell.StyleDefault)
	}

	m.nonVisibleRows = countNoPrintedRows
	return countNoPrintedRows
}

func (m *Menu) DoSelected(t *terminal.Terminal) {
	if m.selectedItem != nil {
		if m.selectedItem.subMenu != nil {
			t.GoToNewView(m.selectedItem.subMenu)
		}
	}
}

type MenuItem struct {
	Title       string
	Description string
	parent      *MenuItem
	child       *MenuItem
	subMenu     *Menu
	isSelected  bool
}

func NewMenuItem(title string) *MenuItem {
	return &MenuItem{
		Title: title,
	}
}

func (mi *MenuItem) WithSubMenu(menu *Menu) *MenuItem {
	mi.subMenu = menu
	return mi
}

func (mi *MenuItem) UnderlineIfSelected() tcell.Style {
	if mi.isSelected {
		return tcell.StyleDefault.Underline(true)
	}
	return tcell.StyleDefault.Underline(false)
}
