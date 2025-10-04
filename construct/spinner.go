package construct

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spinnerModel struct {
	spinner    spinner.Model
	message    string
	subMessage string
	done       bool
	err        error
	quitting   bool
}

type updateMessageMsg string

type taskDoneMsg struct{ err error }

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case updateMessageMsg:
		m.subMessage = string(msg)
		return m, nil

	case taskDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}

	if m.done {
		if m.err != nil {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			return style.Render("✗ " + m.message + " failed\n")
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		return style.Render("✓ " + m.message + " complete\n")
	}

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	msg := spinnerStyle.Render(m.spinner.View()) + " " + m.message
	if m.subMessage != "" {
		subStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		msg += " " + subStyle.Render("("+m.subMessage+")")
	}
	return msg + "\n"
}

// WithSpinner runs a function with a loading spinner
func WithSpinner(message string, fn func() error) error {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := spinnerModel{
		spinner: s,
		message: message,
	}

	p := tea.NewProgram(m)

	// Run the task in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay for spinner to start
		err := fn()
		p.Send(taskDoneMsg{err: err})
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	if m.err != nil {
		return m.err
	}

	return nil
}

// ShowProgress displays a simple progress message
func ShowProgress(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	fmt.Println(style.Render("→ " + message))
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	fmt.Println(style.Render("✓ " + message))
}

// ShowError displays an error message
func ShowError(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	fmt.Println(style.Render("✗ " + message))
}
