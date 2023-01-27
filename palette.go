package palette

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Title string
	Body  string

	position Position
}

func (i *Item) setPosition(horizontal int, vertical int) {
	i.position.h = horizontal
	i.position.v = vertical
}

func (i Item) inFocus(p Position) bool {
	return i.position == p
}

type KeyMap struct {
	CursorLeft  key.Binding
	CursorUp    key.Binding
	CursorRight key.Binding
	CursorDown  key.Binding
	ForceQuit   key.Binding
}

// Postition represents a vertical and a horizontal coordinates of
// and item in the palette matrix.
//
// Coordinates start in the upper left corner of the screen
// with Position{h:0, v:0} being the initial coordinate.
type Position struct {
	h int // horizontal coordinate
	v int // vertical coordinate
}

type Styles struct {
	Item       lipgloss.Style
	ActiveItem lipgloss.Style
}

func defaultStyles() (s Styles) {
	s.Item = lipgloss.NewStyle().Bold(true).
		Margin(2).Padding(1).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))

	s.ActiveItem = s.Item.Copy().
		Background(lipgloss.Color("#94b9f2"))

	return s
}

type Model struct {
	Title string

	width        int
	height       int
	itemWidth    int
	itemsPerLine int
	numLines     int

	styles Styles
	focus  Position
	KeyMap KeyMap
	Help   help.Model

	items []Item
}

func New(items []Item, width, height int) Model {
	maxItemWidth := len(items[0].Title)
	for _, item := range items {
		if len(item.Title) > maxItemWidth {
			maxItemWidth = len(item.Title)
		}
	}

	styles := defaultStyles()
	maxItemWidth += styles.Item.GetHorizontalPadding() + styles.Item.GetHorizontalMargins()

	m := Model{
		Title:     "Palette",
		width:     width,
		height:    height,
		items:     items,
		itemWidth: maxItemWidth,
		focus:     Position{h: 0, v: 0},
		numLines:  1,
		styles:    styles,
		KeyMap:    defaultKeyMap(),
		Help:      help.New(),
	}

	m.calcItemsPerLine()

	return m
}

// calcItemsPerLine calculates a number of items per line
// based on scren and item width
func (m *Model) calcItemsPerLine() {
	m.itemsPerLine = int(math.Max(1, float64(m.width)/float64(m.itemWidth)))
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
		CursorLeft: key.NewBinding(
			key.WithKeys("left", "h"),
		),
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		CursorRight: key.NewBinding(
			key.WithKeys("right", "l"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
		),
	}
}

// SetSize updates various parameters of the model
// related to the size of the screen and number of
// elements displayed on the screen
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.calcItemsPerLine()
	m.numLines = len(m.items) / m.itemsPerLine
}

func (m Model) View() string {
	var (
		output strings.Builder
	)

	if len(m.items) == 0 {
		return "No items found"
	}

	output.WriteString("List of test environments" + "\n\n")

	if len(m.items)%m.itemsPerLine > 0 {
		m.numLines++
	}

	for i := 0; i < m.numLines; i++ {
		itemsRow := ""

		for j := 0; j+i*m.itemsPerLine < len(m.items); j++ {
			curItem := m.items[j+i*m.itemsPerLine]
			curItem.setPosition(j, i) // h: j, v: i

			style := m.styles.Item
			if curItem.inFocus(m.focus) {
				style = m.styles.ActiveItem
			}

			itemsRow = lipgloss.JoinHorizontal(lipgloss.Top, itemsRow,
				style.Width(m.itemWidth).Render(curItem.Title),
			)
		}
		output.WriteString(itemsRow + "\n\n")
	}

	fmt.Fprintf(&output, "\n\nModel width: %d, model height: %d\n", m.width, m.height)

	return output.String()
}

// CursorLeft moves the cursor left
func (m *Model) CursorLeft() {
	if m.focus.h-1 < 0 {
		return
	} else {
		m.focus.h--
	}
}

// CursorUp moves the cursor up
func (m *Model) CursorUp() {
	if m.focus.v-1 < 0 {
		return
	} else {
		m.focus.v--
	}
}

// CursorRight moves the cursor right
func (m *Model) CursorRight() {
	if m.focus.h+1 > m.itemsPerLine {
		return
	} else {
		m.focus.h++
	}
}

// CursorDown moves the cursor down
func (m *Model) CursorDown() {
	if m.focus.v+1 > m.numLines {
		return
	} else {
		m.focus.v++
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.KeyMap.ForceQuit) {
			return m, tea.Quit
		}
	}

	cmds = append(cmds, m.moveCursor(msg))

	return m, tea.Batch(cmds...)
}

// moveCursor handles key pressing events for moving cursor
func (m *Model) moveCursor(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.CursorLeft):
			m.CursorLeft()

		case key.Matches(msg, m.KeyMap.CursorUp):
			m.CursorUp()

		case key.Matches(msg, m.KeyMap.CursorRight):
			m.CursorRight()

		case key.Matches(msg, m.KeyMap.CursorDown):
			m.CursorDown()
		}
	}

	return tea.Batch(cmds...)
}
