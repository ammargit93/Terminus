package tui

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func getAllFiles() []table.Row {
	var files []table.Row
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if slices.Contains(strings.Split(path, "\\"), ".git") {
			return nil
		}
		files = append(files, table.Row{path})
		return nil
	})
	return files
}

type FilePicker struct {
	Table       table.Model
	FileContext []string
}

func InitialiseFilePicker() FilePicker {
	columns := []table.Column{
		{Title: "Files", Width: 65},
	}

	rows := getAllFiles()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false), // start unfocused
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return FilePicker{Table: t}
}

func (m FilePicker) Init() tea.Cmd { return nil }

func (m FilePicker) Update(msg tea.Msg) (FilePicker, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			SelectedRow = m.Table.SelectedRow()[0]
			m.FileContext = append(m.FileContext, SelectedRow)
		}
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m FilePicker) View() string {
	return baseStyle.Render(m.Table.View())
}
