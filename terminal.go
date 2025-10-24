package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ammargit93/terminus/tui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type conversation struct {
	userMessage string
	aiMessage   string
}

// Model defines the main application state.
type model struct {
	chatbox     tui.Chatbox
	keys        tui.KeyMap
	help        help.Model
	modelPicker tui.ModelPickerModel
	lastKey     string
	quitted     bool
	showTable   bool
	viewport    viewport.Model
	LLM         llm
	messages    []conversation
	copyMode    bool
	ready       bool // Add this to track if viewport is ready
}

func newModel() model {
	// Initialize with minimal size, will be updated on first resize
	vp := viewport.New(1, 1)
	vp.KeyMap = viewport.KeyMap{
		PageDown:     key.NewBinding(key.WithKeys("pgdown")),
		PageUp:       key.NewBinding(key.WithKeys("pgup")),
		HalfPageUp:   key.NewBinding(),
		HalfPageDown: key.NewBinding(),
		Up:           key.NewBinding(key.WithKeys("up")),
		Down:         key.NewBinding(key.WithKeys("down")),
	}

	return model{
		chatbox:     tui.NewChatbox(100, 1, 0, "Enter here..."),
		modelPicker: tui.InitialiseModelPicker(),
		keys:        tui.Keys,
		viewport:    vp,
		LLM:         InitialiseModel("llama-3.3-70b-versatile", "groq"),
		help:        help.New(),
		messages:    []conversation{},
		copyMode:    false,
		ready:       false,
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
		if !m.ready {
			// First resize - initialize viewport properly
			m.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.viewport.YPosition = 0
			m.ready = true
			m.updateViewportContent() // Set initial content
		} else {
			// Subsequent resizes
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 10
		}

		m.chatbox.Width = msg.Width - 2
		m.chatbox.Textarea.SetWidth(m.chatbox.Width)
		tui.BottomAlign.Width(m.chatbox.Width)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "up"
			m.viewport.LineUp(1)
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "down"
			m.viewport.LineDown(1)
		case msg.String() == "pgup":
			m.viewport.ViewUp()
		case msg.String() == "pgdown":
			m.viewport.ViewDown()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
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

		case "ctrl+o":
			m.copyMode = !m.copyMode

		case "enter":
			userMessage := strings.TrimSpace(m.chatbox.Textarea.Value())
			var convo conversation

			if userMessage != "" {
				convo.userMessage = userMessage

				aiMessage, err := m.LLM.invoke(userMessage)

				if err != nil {
					aiMessage = "[error: see stderr]"
					fmt.Fprintln(os.Stderr, "DEBUG:", err)
				}
				convo.aiMessage = aiMessage

				// Append conversation
				m.messages = append(m.messages, convo)

				// Update viewport content and scroll to bottom
				m.updateViewportContent()
				m.viewport.GotoBottom()

				m.chatbox.Textarea.SetValue("")
			}

		default:
			if !m.chatbox.Textarea.Focused() {
				cmd = m.chatbox.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	// Update viewport (handles mouse wheel and other viewport-specific inputs)
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update textarea and collect command
	m.chatbox.Textarea, cmd = m.chatbox.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	// Keep table overlay functional
	if m.showTable {
		var tcmd tea.Cmd
		_, tcmd = m.modelPicker.Update(tea.KeyMsg{})
		cmds = append(cmds, tcmd)
	}

	return m, tea.Batch(cmds...)
}

// updateViewportContent updates the viewport with formatted messages
func (m *model) updateViewportContent() {
	if !m.ready {
		return
	}

	var sb strings.Builder
	for _, msg := range m.messages {
		userMsg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#49d765ff")).
			Render("> " + msg.userMessage)

		aiMsg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#3f3fcbff")).
			Render("> " + Wrap(msg.aiMessage, m.chatbox.Width))

		sb.WriteString(userMsg + "\n" + aiMsg + "\n\n")
	}

	m.viewport.SetContent(sb.String())
}

// View renders the UI.
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.showTable {
		return m.renderWithOverlay()
	}
	return m.renderBase()
}

func (m model) renderBase() string {
	// Use the viewport for rendering content
	content := tui.TerminusStyle.Render(tui.Terminus) + "\n\n" +
		m.viewport.View() + "\n" +
		m.renderInputHelp()

	return content
}

func (m model) renderInputHelp() string {
	return m.chatbox.Textarea.View() + m.help.View(m.keys)
}

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

func main() {
	p := tea.NewProgram(newModel(), tea.WithMouseCellMotion())

	// Run program and capture the returned model
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	m := finalModel.(model)

	for _, msg := range m.messages {
		fmt.Printf("> %s\n> %s\n\n", msg.userMessage, msg.aiMessage)
	}

}
