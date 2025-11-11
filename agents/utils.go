package agents

import (
	"encoding/json"
	"os"
	"strings"
)

type Argument struct {
	ActionName string   `json:"actionName"`
	ArgNames   []string `json:"argNames"`
}
type ResponseMessage struct {
	Action  []string   `json:"action"`
	Message string     `json:"message"`
	Code    string     `json:"code"`
	Args    []Argument `json:"args"`
	ToolUse bool       `json:"tooluse"`
}

func ParseJSON(content string) ResponseMessage {
	var v ResponseMessage
	json.Unmarshal([]byte(content), &v)
	return v
}

func getCWD(path string) string {
	cwd, _ := os.Getwd()
	finalPath := cwd + "\\" + path
	finalPath = strings.ReplaceAll(finalPath, "/", "\\")
	finalPath = strings.ReplaceAll(finalPath, `\\`, `\`)
	return finalPath
}

var toolKeywords = []string{
	"write", "create", "make", "save", "generate", "export", "delete", "remove", "update",
	"read", "open", "download", "upload", "print", "send", "execute", "build",
}

// Simple function to classify
func IsToolPrompt(prompt string) bool {
	prompt = strings.ToLower(prompt)
	for _, kw := range toolKeywords {
		if strings.Contains(prompt, kw) {
			return true
		}
	}
	return false
}
