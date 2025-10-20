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
	TerminusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#59db77ff"))

	BgStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#131824")). // blue background
		Foreground(lipgloss.Color("0"))        // white text\

	BottomAlign = lipgloss.NewStyle().Align(lipgloss.Bottom, lipgloss.Center)
)
