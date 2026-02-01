package tui

import (
	"ensemble/internal/util"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		PaddingLeft(2)

	emptyStateStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(2).
		Italic(true)

	pickerContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("62"))
)

type ProjectLoader struct {
	filepicker     filepicker.Model
	loadedProjects []string
	windowHeight   int
	windowWidth    int
}

func (p *ProjectLoader) Init() tea.Cmd {
	return p.filepicker.Init()
}

func (p *ProjectLoader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case UpMsg:
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case DownMsg:
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case EnterMsg:
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case tea.WindowSizeMsg:
		p.windowHeight = msgType.Height
		p.windowWidth = msgType.Width
	}

	var cmd tea.Cmd
	p.filepicker, cmd = p.filepicker.Update(msg)

	if didSelect, path := p.filepicker.DidSelectFile(msg); didSelect {
		if loaded := util.StartProject(path); loaded {
			p.loadedProjects = append(p.loadedProjects, path)
		}
	}

	return p, cmd
}

func (p *ProjectLoader) View() string {
	p.filepicker.SetHeight(p.calculatePickerHeight())

	var s strings.Builder

	s.WriteString(titleStyle.Render("Loaded Projects") + "\n")

	if len(p.loadedProjects) == 0 {
		s.WriteString(emptyStateStyle.Render("No projects loaded.") + "\n")
	} else {
		for _, proj := range p.loadedProjects {
			s.WriteString(itemStyle.Render("✓ "+proj) + "\n")
		}
	}

	s.WriteString("\n")
	pickerContent := "Select a File:\n" + p.filepicker.View()
	s.WriteString(pickerContainerStyle.Render(pickerContent))

	return s.String()
}

func (p *ProjectLoader) calculatePickerHeight() int {
	projectLines := len(p.loadedProjects)
	if projectLines == 0 {
		projectLines = 1
	}

	usedLines := 7 + projectLines
	return max(p.windowHeight-usedLines, 5)
}

func CreateProjectLoader() *ProjectLoader {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".json"}
	fp.ShowHidden = true

	if dir, ok := os.LookupEnv("UPLOAD_DIR"); ok {
		fp.CurrentDirectory = dir
	} else {
		fp.CurrentDirectory, _ = os.Getwd()
	}

	return &ProjectLoader{
		filepicker:     fp,
		loadedProjects: util.GetLoadedProjects(),
	}
}
