package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// custom messages
type navMsg string
type UpMsg struct{}
type DownMsg struct{}
type EnterMsg struct{}

type Tui struct {
	screens map[string]tea.Model
	current string
}

func (t Tui) Init() tea.Cmd {
	if s, ok := t.screens[t.current]; ok {
		return s.Init()
	}
	return nil
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "up", "k":
			msg = UpMsg{}
		case "down", "j":
			msg = DownMsg{}
		case "enter":
			msg = EnterMsg{}
		}

	case navMsg:
		target := string(m)
		if s, ok := t.screens[target]; ok {
			t.current = target
			return t, s.Init()
		}
		return t, nil
	}

	if active, ok := t.screens[t.current]; ok {
		newModel, newCmd := active.Update(msg)
		t.screens[t.current] = newModel
		return t, newCmd
	}

	return t, nil
}

func (t Tui) View() string {
	if s, ok := t.screens[t.current]; ok {
		return s.View()
	}
	return "Screen not found."
}

func CreateTui() Tui {
	return Tui{
		screens: map[string]tea.Model{
			"main_menu": CreateMainMenu(),
		},
		current: "main_menu",
	}
}
