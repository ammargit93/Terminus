package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ammargit93/terminus/tui"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	LLM         llm
	messages    []string
}

func newModel() model {
	return model{
		chatbox:     tui.NewChatbox(100, 1, 0, "Enter here..."),
		modelPicker: tui.InitialiseModelPicker(),
		keys:        tui.Keys,
		LLM:         InitialiseModel("llama-3.3-70b-versatile", "groq"),
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
				resp, err := m.LLM.invoke(text)
				if err != nil {
					fmt.Fprintln(os.Stderr, "DEBUG:", err)
				} else {
					fmt.Fprintln(os.Stdout, resp)
				}
			}

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
	base := m.renderBase()
	if !m.showTable {
		return base + m.renderInputHelp()
	}
	return m.renderWithOverlay()
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
