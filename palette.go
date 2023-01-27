package palette

import (
	"fmt"
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

func (i *Item) setPosition(row int, col int) {
	i.position.row = row
	i.position.col = col
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
// with Position{row:0, col:0} being the initial coordinate.
type Position struct {
	row int
	col int
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

type State int

const (
	initializing State = iota
	ready
)

type Model struct {
	Title string

	width        int
	height       int
	itemWidth    int
	itemsPerLine int
	numLines     int

	state  State
	styles Styles
	focus  Position
	KeyMap KeyMap
	Help   help.Model

	items []Item
}

func (m *Model) StateReady() {
	m.state = ready
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
		focus:     Position{row: 0, col: 0},
		styles:    styles,
		KeyMap:    defaultKeyMap(),
		Help:      help.New(),
	}

	m.calcItemsPerLine()
	m.calcNumberOfLines()

	return m
}

// calcItemsPerLine calculates the number of items per line
// based on the screen and item width
func (m *Model) calcItemsPerLine() {
	ipl := m.width / m.itemWidth
	if m.width%m.itemWidth > 0 {
		ipl--
	}

	m.itemsPerLine = max(1, ipl)
}

// calcNumberOfLines calculates a number of lines on the screen
func (m *Model) calcNumberOfLines() {
	m.numLines = len(m.items) / m.itemsPerLine
	if len(m.items)%m.itemsPerLine > 0 {
		m.numLines++
	}
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
	m.calcNumberOfLines()
}

func (m Model) View() string {
	var (
		output strings.Builder
	)

	if m.state == initializing {
		return "Initializing..."
	}

	if len(m.items) == 0 {
		return "No items found"
	}

	output.WriteString("List of test environments" + "\n\n")

	for lineNumber := 0; lineNumber < m.numLines; lineNumber++ {
		itemsRow := ""

		start := lineNumber * m.itemsPerLine
		end := start + m.itemsPerLine
		if end > len(m.items) {
			end = len(m.items)
		}

		for columnNumber, item := range m.items[start:end] {
			item.setPosition(lineNumber, columnNumber)

			style := m.styles.Item
			if item.inFocus(m.focus) {
				style = m.styles.ActiveItem
			}

			itemsRow = lipgloss.JoinHorizontal(lipgloss.Top, itemsRow,
				style.Width(m.itemWidth).Render(item.Title),
			)
		}

		output.WriteString(itemsRow + "\n\n")
	}

	fmt.Fprintf(&output, "\n\nModel width: %d, model height: %d\n", m.width, m.height)

	return output.String()
}

// CursorLeft moves the cursor left
func (m *Model) CursorLeft() {
	if m.focus.col > 0 {
		m.focus.col--
	}
}

// CursorUp moves the cursor up
func (m *Model) CursorUp() {
	if m.focus.row > 0 {
		m.focus.row--
	}
}

// CursorRight moves the cursor right
func (m *Model) CursorRight() {
	m.focus.col++

	if m.focus.col == m.itemsPerLine {
		m.focus.col--
		return
	}

	if m.focus.row == m.numLines-1 && len(m.items)%m.itemsPerLine > 0 {
		if len(m.items)%m.itemsPerLine-1 < m.focus.col {
			m.focus.col--
		}
	}
}

// CursorDown moves the cursor down
func (m *Model) CursorDown() {
	m.focus.row++

	if m.focus.row == m.numLines {
		m.focus.row--
		return
	}

	if m.focus.row == m.numLines-1 && len(m.items)%m.itemsPerLine > 0 {
		if len(m.items)%m.itemsPerLine-1 < m.focus.col {
			m.focus.row--
		}
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
