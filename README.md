# lip-palette

Terminal UI component that allows you creating a palette (matrix) of items. Based on [lipglos](https://github.com/charmbracelet/lipgloss)

# Usage

Get the lip-palette module

```bash
go get github.com/misha-slyusarev/lip-palette@v0.1.0
```

Then in your code use it together with [bubbletea](https://github.com/charmbracelet/bubbletea)
```go
package ui

import (
  tea "github.com/charmbracelet/bubbletea"
  palette "github.com/misha-slyusarev/lip-palette"
)

type mainPageModel struct {
  plt palette.Model
}

func (m mainPageModel) Init() tea.Cmd {
  return tea.EnterAltScreen
}

...
```
