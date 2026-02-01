package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Custom Messages
type navMsg string
type UpMsg struct{}
type DownMsg struct{}
type EnterMsg struct{}
type cancelMsg struct{}
type Stack[T any] struct {
	s []T
}

func (st *Stack[T]) Push(v T) {
	st.s = append(st.s, v)
}

func (st *Stack[T]) Pop() (T, bool) {
	if len(st.s) == 0 {
		var zero T
		return zero, false
	}

	l := len(st.s)
	res := st.s[l-1]
	st.s = st.s[:l-1]
	return res, true
}

func (st *Stack[T]) Peek() (T, bool) {
	if len(st.s) == 0 {
		var zero T
		return zero, false
	}
	return st.s[len(st.s)-1], true
}

type Tui struct {
	screens map[string]tea.Model
	stack   Stack[string]
	root    string
	width   int
	height  int
}

func (t *Tui) Init() tea.Cmd {
	if s, ok := t.screens[t.root]; ok {
		return s.Init()
	}
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = m.Width
		t.height = m.Height

	case tea.KeyMsg:
		switch m.String() {
		case "ctrl+c", "q":
			return t, tea.Quit

		case "esc", "b":
			t.stack.Pop()
			screen, exists := t.stack.Peek()
			if !exists {
				return t, tea.Quit
			}

			if s, ok := t.screens[screen]; ok {
				updatedModel, cmd := s.Update(tea.WindowSizeMsg{Width: t.width, Height: t.height})
				t.screens[screen] = updatedModel
				return t, cmd
			}
			return t, nil

		case "up", "k":
			msg = UpMsg{}
		case "down", "j":
			msg = DownMsg{}
		case "enter", "u":
			msg = EnterMsg{}
		case "c":
			msg = cancelMsg{}
		}

	case navMsg:
		target := string(m)
		if s, ok := t.screens[target]; ok {
			t.stack.Push(target)
			updatedModel, cmd := s.Update(tea.WindowSizeMsg{Width: t.width, Height: t.height})
			t.screens[target] = updatedModel

			return t, tea.Batch(cmd, s.Init())
		}
		return t, nil
	}

	screen, found := t.stack.Peek()
	if !found {
		screen = t.root
		t.stack.Push(screen)
	}

	if active, ok := t.screens[screen]; ok {
		newModel, newCmd := active.Update(msg)
		t.screens[screen] = newModel
		return t, newCmd
	}

	return t, nil
}

func (t *Tui) View() string {
	screen, found := t.stack.Peek()
	if !found {
		return "Error: No active screen found."
	}

	if s, ok := t.screens[screen]; ok {
		return s.View()
	}
	return "Error: Screen model missing."
}

func CreateTui() *Tui {
	st := Stack[string]{}
	st.Push("main_menu")

	return &Tui{
		screens: map[string]tea.Model{
			"main_menu":        CreateMainMenu(),
			"Load New Project": CreateProjectLoader(),
		},
		stack: st,
		root:  "main_menu",
	}
}
