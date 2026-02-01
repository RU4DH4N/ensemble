package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	titleStyle   = lipgloss.NewStyle().Margin(1, 0, 1, 2).Padding(0, 1).Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
)

type MainMenu struct {
	choices []string
	cursor  int
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case UpMsg:
		if m.cursor > 0 {
			m.cursor--
		}
	case DownMsg:
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case EnterMsg:
		target := m.choices[m.cursor]

		return m, func() tea.Msg {
			return navMsg(target)
		}
	}

	return m, nil
}

func (m MainMenu) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(" Select an Option "))
	s.WriteString("\n\n")

	for i, choice := range m.choices {
		cursor := "  "
		label := choice
		if m.cursor == i {
			cursor = focusedStyle.Render("> ")
			label = focusedStyle.Render(choice)
		}

		s.WriteString(fmt.Sprintf("%s %s\n", cursor, label))
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" (j/k or arrows to move, enter to select, q to quit) "))
	s.WriteString("\n")

	return s.String()
}

func CreateMainMenu() MainMenu {
	return MainMenu{
		choices: []string{"Current Projects Overview", "Load New Project", "Settings"},
	}
}
