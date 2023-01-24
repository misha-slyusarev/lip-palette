package palette

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Title string
	Body  string
}

type Model struct {
	Title string

	width        int
	height       int
	itemWidth    int
	itemsPerLine int
	cursor       int
	Help         help.Model

	// The master set of items we're working with.
	items []Item
}

func New(items []Item, width, height int) Model {
	maxItemWidth := len(items[0].Title)
	for _, item := range items {
		if len(item.Title) > maxItemWidth {
			maxItemWidth = len(item.Title)
		}
	}

	m := Model{
		Title:        "List",
		width:        width,
		height:       height,
		items:        items,
		itemWidth:    maxItemWidth,
		itemsPerLine: int(math.Max(1, float64(width)/float64(maxItemWidth))),
		Help:         help.New(),
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.itemsPerLine = width / m.itemWidth
}

var itemStyle = lipgloss.NewStyle().Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

func (m Model) View() string {
	var (
		output strings.Builder
	)

	if len(m.items) == 0 {
		return "No items found"
	}

	output.WriteString("List of test environments" + "\n\n")

	itemsRow := ""

	for i, item := range m.items {
		itemsRow = lipgloss.JoinHorizontal(lipgloss.Top, itemsRow,
			itemStyle.Copy().Width(m.itemWidth+4).Margin(2).Padding(1).Render(item.Title),
		)

		if (i+1)%m.itemsPerLine == 0 {
			output.WriteString(itemsRow + "\n\n")

			itemsRow = ""
		}
	}

	if m.itemsPerLine > len(m.items) || itemsRow != "" {
		output.WriteString(itemsRow + "\n\n")
	}

	fmt.Fprintf(&output, "\n\nModel width: %d, model height: %d\n", m.width, m.height)

	return output.String()
}
