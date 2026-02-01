package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stateFilePicker = iota
	stateFilePreview
	stateUploading
)

type clearErrorMsg struct{}
type FileUploader struct {
	filepicker     filepicker.Model
	viewport       viewport.Model
	progressBar    progress.Model
	selectedFile   string
	fileContent    string
	uploadProgress float64
	windowHeight   int
	windowWidth    int
	state          int
	err            error
}

func (f *FileUploader) Init() tea.Cmd {
	f.viewport = viewport.New(80, 20)
	f.progressBar = progress.New(progress.WithDefaultGradient())
	return f.filepicker.Init()
}

func (f *FileUploader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case UpMsg:
		if f.state == stateFilePreview {
			f.viewport.LineUp(1)
			return f, nil
		}
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case DownMsg:
		if f.state == stateFilePreview {
			f.viewport.LineDown(1)
			return f, nil
		}
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case EnterMsg:
		if f.state == stateFilePicker {
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		}
	case tea.WindowSizeMsg:
		f.windowHeight = msgType.Height
		f.windowWidth = msgType.Width
		f.filepicker.SetHeight(msgType.Height - 4)

		// Set viewport size for file preview (reserve space for header and buttons)
		viewportHeight := msgType.Height - 10
		if viewportHeight < 10 {
			viewportHeight = 10
		}
		f.viewport.Width = msgType.Width - 4
		f.viewport.Height = viewportHeight
	case clearErrorMsg:
		f.err = nil
	case cancelMsg:
		if f.state == stateFilePreview {
			f.state = stateFilePicker
			f.selectedFile = ""
			f.fileContent = ""
			f.viewport.SetContent("")
		}
	}

	if f.state == stateFilePicker {
		var cmd tea.Cmd
		f.filepicker, cmd = f.filepicker.Update(msg)

		if didSelect, path := f.filepicker.DidSelectFile(msg); didSelect {
			f.selectedFile = path

			content, err := os.ReadFile(path)
			if err != nil {
				f.err = errors.New("Failed to read file: " + err.Error())
				f.selectedFile = ""
				return f, tea.Batch(cmd, clearErrorAfter(2*time.Second))
			}

			f.fileContent = string(content)
			f.viewport.SetContent(f.fileContent)
			f.viewport.GotoTop()
			f.state = stateFilePreview
		}

		return f, cmd
	}

	return f, nil
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (f *FileUploader) View() string {
	var s strings.Builder

	switch f.state {
	case stateFilePicker:
		s.WriteString("\n  ")
		if f.err != nil {
			s.WriteString(f.filepicker.Styles.DisabledFile.Render(f.err.Error()))
		} else {
			s.WriteString("Pick a file (Press 'esc' or 'q' to exit):")
		}
		s.WriteString("\n\n" + f.filepicker.View() + "\n")

	case stateFilePreview:
		s.WriteString("\n  ")
		s.WriteString("Selected file: " + f.filepicker.Styles.Selected.Render(f.selectedFile))
		s.WriteString("\n\n")
		s.WriteString("File content (use arrow keys to scroll):\n")
		s.WriteString("─────────────────────────────────────────────────\n")

		// Use viewport to display file content
		s.WriteString(f.viewport.View())

		s.WriteString("\n─────────────────────────────────────────────────\n")

		// Show scroll position indicator
		scrollPercent := int(f.viewport.ScrollPercent() * 100)
		scrollInfo := ""
		if scrollPercent >= 100 {
			scrollInfo = "(bottom)"
		} else if scrollPercent <= 0 {
			scrollInfo = "(top)"
		} else {
			scrollInfo = fmt.Sprintf("(%d%%)", scrollPercent)
		}

		s.WriteString("  ↑/↓ " + scrollInfo + "\n\n")
		s.WriteString("  [c] Cancel  |  [u] Upload\n")

	case stateUploading:
		s.WriteString("\n  ")
		s.WriteString("Uploading file...\n")
		s.WriteString("\n")

		// Use progress bar
		s.WriteString("  " + f.progressBar.ViewAs(f.uploadProgress) + "\n")
		s.WriteString("\n")
	}

	return s.String()
}

func CreateFileUploader() *FileUploader {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".json"}
	if dir, err := os.UserHomeDir(); err == nil {
		fp.CurrentDirectory = dir
	} else {
		fp.CurrentDirectory = "."
	}

	vp := viewport.New(80, 20)
	pb := progress.New(progress.WithDefaultGradient())

	return &FileUploader{
		filepicker:  fp,
		viewport:    vp,
		progressBar: pb,
		state:       stateFilePicker,
	}
}
