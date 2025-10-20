package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
)

type Chatbox struct {
	Textarea textarea.Model
	Width    int
	Height   int
	Padding  int
}

func NewChatbox(width, height, padding int, placeholder string) Chatbox {
	ta := textarea.New()
	ta.SetWidth(width)
	ta.SetHeight(height)
	ta.ShowLineNumbers = false
	ta.Placeholder = placeholder
	ta.Prompt = "> "

	return Chatbox{
		Textarea: ta,
		Width:    width,
		Height:   height,
		Padding:  padding,
	}
}
