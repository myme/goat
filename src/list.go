package goat

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list   list.Model
	choice *Item
	err    error
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.choice = nil
			m.err = fmt.Errorf("aborted")
			return m, tea.Quit
		case "enter":
			item, ok := m.list.SelectedItem().(Item)
			if !ok {
				m.err = fmt.Errorf("invalid selection")
				return m, tea.Quit
			}
			m.choice = &item
			m.err = nil
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

type Item struct {
	Index  int
	Text   string
	Desc   string
	Filter string
}

func (i Item) Title() string       { return i.Text }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Filter }

// Displays a pretty selection list of [Item] and return the selected [Item].
func SelectFromList(title string, items []Item) (*Item, error) {
	listItems := make([]list.Item, len(items))
	for i := range items {
		listItems[i] = items[i]
	}
	m := model{list: list.New(listItems, list.NewDefaultDelegate(), 0, 10)}
	m.list.Title = title

	res, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return nil, err
	}
	if res.(model).err != nil {
		return nil, res.(model).err
	}

	return res.(model).choice, nil
}
