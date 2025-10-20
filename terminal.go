package main

import (
	"fmt"
	"os"
	"practice/tui"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model

type model struct {
	chatbox tui.Chatbox
	Keys    tui.KeyMap
	help    help.Model
	lastKey string
	quitted bool
}

func newModel() model {
	chatbox := tui.NewChatbox(100, 1, 0, "Enter here...") // helper constructor
	return model{
		chatbox: chatbox,
		Keys:    tui.Keys,
		help:    help.New(),
		quitted: false,
	}
}

func (m model) Init() tea.Cmd {
	return m.chatbox.Textarea.Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Dynamically set textarea width based on terminal width
		m.chatbox.Width = msg.Width - 2 // leave some padding
		m.chatbox.Height = msg.Height
		m.chatbox.Textarea.SetWidth(m.chatbox.Width)
		tui.BottomAlign.Width(m.chatbox.Width)

		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			m.lastKey = "up"
		case key.Matches(msg, m.Keys.Down):
			m.lastKey = "down"
		case key.Matches(msg, m.Keys.Left):
			m.lastKey = "left"
		case key.Matches(msg, m.Keys.Right):
			m.lastKey = "right"

		case key.Matches(msg, m.Keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.Keys.Quit):
			m.quitted = true
			return m, tea.Quit
		}

		switch msg.String() {
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

func (m model) View() string {
	input := m.chatbox.Textarea.View()
	helpView := m.help.View(m.Keys)

	availableHeight := m.chatbox.Height
	inputHeight := lipgloss.Height(input)

	// space below input (offset from bottom)
	offset := 2 // number of lines you want "above bottom"

	padding := availableHeight - inputHeight - offset - 6
	if padding < 0 {
		padding = 0
	}

	output := tui.TerminusStyle.Render(tui.Terminus) + strings.Repeat("\n", padding) + input + helpView

	return output
}

func main() {

	_, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
