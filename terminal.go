package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ammargit93/terminus/tui"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model defines the main application state.
type model struct {
	chatbox     tui.Chatbox
	keys        tui.KeyMap
	help        help.Model
	modelPicker tui.ModelPickerModel
	lastKey     string
	quitted     bool
	showTable   bool
	messages    []string
}

func newModel() model {
	return model{
		chatbox:     tui.NewChatbox(100, 1, 0, "Enter here..."),
		modelPicker: tui.InitialiseModelPicker(),
		keys:        tui.Keys,
		help:        help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.chatbox.Textarea.Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// Resize chatbox dynamically
		m.chatbox.Width = msg.Width - 2
		m.chatbox.Height = msg.Height
		m.chatbox.Textarea.SetWidth(m.chatbox.Width)
		tui.BottomAlign.Width(m.chatbox.Width)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "up"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "down"
		case key.Matches(msg, m.keys.Left):
			m.lastKey = "left"
		case key.Matches(msg, m.keys.Right):
			m.lastKey = "right"

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitted = true
			return m, tea.Quit
		}
		switch msg.String() {
		case "tab":
			if m.chatbox.Textarea.Focused() {
				m.chatbox.Textarea.Blur()
			} else {
				m.chatbox.Textarea.Focus()
			}
		case "ctrl+c":
			m.quitted = true
			return m, tea.Quit
		case "esc":
			m.showTable = !m.showTable
		case "enter":
			text := strings.TrimSpace(m.chatbox.Textarea.Value())
			if text != "" {
				m.messages = append(m.messages, text)
				m.chatbox.Textarea.SetValue("") // clear after sending
			}

		default:
			if !m.chatbox.Textarea.Focused() {
				cmd = m.chatbox.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.chatbox.Textarea, cmd = m.chatbox.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	// If the table is shown, route input to it as well so the table can react to keys
	if m.showTable {
		var tcmd tea.Cmd
		_, tcmd = m.modelPicker.Update(tea.KeyMsg{}) // harmless placeholder so table styles stay consistent
		cmds = append(cmds, tcmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI.
func (m model) View() string {
	input := m.chatbox.Textarea.View()
	helpView := m.help.View(m.keys)

	// Build chat history string
	history := ""
	for _, msg := range m.messages {
		history += lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#87CEEB")). // light blue text
			Render("> "+msg) + "\n"
	}

	availableHeight := m.chatbox.Height
	inputHeight := lipgloss.Height(input)
	historyHeight := lipgloss.Height(history)

	offset := 3
	padding := availableHeight - inputHeight - historyHeight - offset - 6
	if padding < 0 {
		padding = 0
	}

	// base holds the unchanged main screen (logo + history + space where overlay can go + input/help)
	base := tui.TerminusStyle.Render(tui.Terminus) +
		"\n\n" + history +
		strings.Repeat("\n", padding)

	// when table is NOT shown, just render base + input + help
	if !m.showTable {
		return base + input + helpView
	}

	// ---- TABLE OVERLAY: render centered and with a fixed max width ----
	overlayWidth := m.chatbox.Width - 8
	if overlayWidth < 40 {
		overlayWidth = m.chatbox.Width
	}

	tableContent := m.modelPicker.View()

	tableBox := lipgloss.NewStyle().
		Width(overlayWidth).
		Padding(0, 1).
		Align(0.25, lipgloss.Top). // top inside its allocated area
		Render(tableContent)

	return tui.TerminusStyle.Render(tui.Terminus) + strings.Repeat("\n", 14) + tableBox + "\n" + input + helpView

}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
