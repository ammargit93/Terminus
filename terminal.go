package main

import (
	"fmt"
	"os"
	"practice/tui"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model

type model struct {
	chatbox tui.Chatbox
	quitted bool
}

func newModel() model {
	chatbox := tui.NewChatbox(100, 1, 0, "Enter here...") // helper constructor
	return model{chatbox: chatbox}
}

func (m model) Init() tea.Cmd {
	return m.chatbox.Textarea.Focus()
}

var BottomAlign = lipgloss.NewStyle().Align(lipgloss.Bottom, lipgloss.Center)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Dynamically set textarea width based on terminal width
		m.chatbox.Width = msg.Width - 2 // leave some padding
		m.chatbox.Height = msg.Height
		m.chatbox.Textarea.SetWidth(m.chatbox.Width)
		BottomAlign.Width(m.chatbox.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitted = true
			return m, tea.Quit
		case "tab":
			m.chatbox.Textarea.Blur()
		default:
			if !m.chatbox.Textarea.Focused() {

				cmd = m.chatbox.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.chatbox.Textarea, cmd = m.chatbox.Textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

var bgStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#131824")). // blue background
	Foreground(lipgloss.Color("0"))        // white text

var terminusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#159f35ff"))

func (m model) View() string {
	input := m.chatbox.Textarea.View()

	availableHeight := m.chatbox.Height
	inputHeight := lipgloss.Height(input)

	// space below input (offset from bottom)
	offset := 2 // number of lines you want "above bottom"

	padding := availableHeight - inputHeight - offset - 6
	if padding < 0 {
		padding = 0
	}

	output := terminusStyle.Render(tui.Terminus) + strings.Repeat("\n", padding) + input

	return bgStyle.Render(output)
}

func main() {

	_, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
