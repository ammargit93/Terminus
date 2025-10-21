package tui

import "github.com/charmbracelet/lipgloss"

var Terminus = `
████████╗███████╗██████╗ ███╗   ███╗██╗███╗   ██╗██╗   ██╗███████╗
╚══██╔══╝██╔════╝██╔══██╗████╗ ████║██║████╗  ██║██║   ██║██╔════╝
   ██║   █████╗  ██████╔╝██╔████╔██║██║██╔██╗ ██║██║   ██║███████╗
   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║██║██║╚██╗██║██║   ██║╚════██║
   ██║   ███████╗██║  ██║██║ ╚═╝ ██║██║██║ ╚████║╚██████╔╝███████║
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚══════╝                                                                 
`
var (
	TerminusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#59db77ff")) // leafgreen

	BgStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#131824")). // blue background
		Foreground(lipgloss.Color("0"))        // white text\

	BottomAlign = lipgloss.NewStyle().Align(lipgloss.Bottom, lipgloss.Center)

	ModelPickerAlign = lipgloss.NewStyle().Align(0.75, 0.3)

	BaseTableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)
