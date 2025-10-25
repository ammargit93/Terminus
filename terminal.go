package main

import (
	"strings"

	"github.com/ammargit93/terminus/tui"
	"github.com/ammargit93/terminus/vector"
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

type model struct {
	chatbox     tui.Chatbox
	keys        tui.KeyMap
	help        help.Model
	filePicker  tui.FilePicker
	showTable   bool
	viewport    viewport.Model
	fileContext []string
	LLM         llm
	embedding   vector.Embedding
	messages    []conversation
	copyMode    bool
	ready       bool
}

func newModel() model {
	vp := viewport.New(1, 1)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Up:       key.NewBinding(key.WithKeys("up")),
		Down:     key.NewBinding(key.WithKeys("down")),
	}

	return model{
		chatbox:    tui.NewChatbox(100, 1, 0, "Enter here..."),
		filePicker: tui.InitialiseFilePicker(),
		embedding:  vector.InitialiseEmbeddingModel(),
		keys:       tui.Keys,
		viewport:   vp,
		LLM:        InitialiseModel("llama-3.3-70b-versatile", "groq"),
		help:       help.New(),
		messages:   []conversation{},
		copyMode:   false,
		ready:      false,
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
			m.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.ready = true
			m.updateViewportContent()
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 10
		}
		m.chatbox.Width = msg.Width - 2
		m.chatbox.Textarea.SetWidth(m.chatbox.Width)
		tui.BottomAlign.Width(m.chatbox.Width)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		if m.showTable {
			// Send keys to table only
			var tcmd tea.Cmd
			m.filePicker, tcmd = m.filePicker.Update(msg)
			cmds = append(cmds, tcmd)
			if msg.String() == "esc" {
				m.showTable = false
				m.fileContext = m.filePicker.FileContext
				vector.CallCohere(m.fileContext)
				m.filePicker.FileContext = []string{}

				m.chatbox.Textarea.Focus()
			}
			return m, tea.Batch(cmds...)
		}

		// Normal chatbox + viewport keys
		switch {
		case key.Matches(msg, m.keys.Up):
			m.viewport.LineUp(1)
		case key.Matches(msg, m.keys.Down):
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
			return m, tea.Quit
		case "esc":
			m.showTable = true
			m.filePicker.Table.Focus()
		case "ctrl+o":
			m.copyMode = !m.copyMode
		case "enter":
			userMessage := strings.TrimSpace(m.chatbox.Textarea.Value())
			if userMessage != "" {
				aiMsg, err := m.LLM.invoke(userMessage)
				if err != nil {
					aiMsg = "[error]"
				}
				m.messages = append(m.messages, conversation{userMessage, aiMsg})
				m.updateViewportContent()
				m.viewport.GotoBottom()
				m.chatbox.Textarea.SetValue("")
			}
		}
	}

	// Update viewport
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update textarea
	m.chatbox.Textarea, cmd = m.chatbox.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) updateViewportContent() {
	if !m.ready {
		return
	}
	var sb strings.Builder
	for _, msg := range m.messages {
		userMsg := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#49d765ff")).
			Render("> " + msg.userMessage)
		aiMsg := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3f3fcbff")).
			Render("> " + Wrap(msg.aiMessage, m.chatbox.Width))
		sb.WriteString(userMsg + "\n" + aiMsg + "\n\n")
	}
	m.viewport.SetContent(sb.String())
}

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
	return tui.TerminusStyle.Render(tui.Terminus) + "\n\n" +
		m.viewport.View() + "\n" +
		m.chatbox.Textarea.View() + m.help.View(m.keys)
}

func (m model) renderWithOverlay() string {
	width := m.chatbox.Width - 8
	if width < 40 {
		width = m.chatbox.Width
	}
	tableBox := lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Align(0.25, lipgloss.Top).
		Render(m.filePicker.View())
	return tui.TerminusStyle.Render(tui.Terminus) +
		strings.Repeat("\n", 14) + tableBox + "\n" +
		m.chatbox.Textarea.View() + m.help.View(m.keys)
}

// func main() {
// 	p := tea.NewProgram(newModel(), tea.WithMouseCellMotion())
// 	finalModel, err := p.Run()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		os.Exit(1)
// 	}

// 	m := finalModel.(model)
// 	for _, msg := range m.messages {
// 		fmt.Printf("> %s\n> %s\n\n", msg.userMessage, msg.aiMessage)
// 	}
// }
