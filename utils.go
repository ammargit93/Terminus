package main

import (
	"strings"

	"github.com/ammargit93/terminus/tui"
	"github.com/charmbracelet/lipgloss"
)

// renderBase renders logo + chat history + padding
func (m model) renderBase() string {
	history := m.renderHistory()

	availableHeight := m.chatbox.Height
	inputHeight := lipgloss.Height(m.chatbox.Textarea.View())
	historyHeight := lipgloss.Height(history)

	offset := 3
	padding := availableHeight - inputHeight - historyHeight - offset - 6
	if padding < 0 {
		padding = 0
	}

	return tui.TerminusStyle.Render(tui.Terminus) +
		"\n\n" + history +
		strings.Repeat("\n", padding)
}

// renderHistory renders chat messages
func (m model) renderHistory() string {
	var sb strings.Builder
	for _, msg := range m.messages {
		sb.WriteString(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#87CEEB")).
				Render("> "+msg) + "\n",
		)
	}
	return sb.String()
}

// renderInputHelp renders the input box and help
func (m model) renderInputHelp() string {
	return m.chatbox.Textarea.View() + m.help.View(m.keys)
}

// renderWithOverlay handles the model picker overlay
func (m model) renderWithOverlay() string {
	overlayWidth := m.chatbox.Width - 8
	if overlayWidth < 40 {
		overlayWidth = m.chatbox.Width
	}

	tableBox := lipgloss.NewStyle().
		Width(overlayWidth).
		Padding(0, 1).
		Align(0.25, lipgloss.Top).
		Render(m.modelPicker.View())

	return tui.TerminusStyle.Render(tui.Terminus) +
		strings.Repeat("\n", 14) + tableBox +
		"\n" + m.renderInputHelp()
}
