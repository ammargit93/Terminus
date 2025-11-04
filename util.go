package main

import (
	"encoding/json"
	"strings"
	"unicode"
)

// Wrap wraps text at the given width while preserving bullets, indents, and spacing.
func Wrap(content string, width int) string {
	if width <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var wrapped []string

	for _, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if len(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}
		prefix := ""
		if len(line) > 1 && (strings.HasPrefix(line, "* ") ||
			strings.HasPrefix(line, "- ") ||
			strings.HasPrefix(line, "> ")) {
			prefix = line[:2]
			line = line[2:]
		} else if len(line) > 0 && (line[0] == '\t' || line[0] == ' ') {
			for i, r := range line {
				if r != ' ' && r != '\t' {
					prefix = line[:i]
					line = line[i:]
					break
				}
			}
		}

		words := strings.Fields(line)
		currLine := prefix
		currLen := len(currLine)

		for _, word := range words {
			if currLen+len(word)+1 > width {
				wrapped = append(wrapped, currLine)
				currLine = prefix + word
				currLen = len(currLine)
			} else {
				if currLen > len(prefix) {
					currLine += " "
					currLen++
				}
				currLine += word
				currLen += len(word)
			}
		}
		wrapped = append(wrapped, currLine)
	}
	ret := boldenContent(strings.Join(wrapped, "\n"))
	return ret
}

func boldenContent(content string) string {
	var result string = ""

	var i int = 0

	for i < len(content) {
		if i+1 < len(content) && content[i] == '*' && content[i+1] == '*' {

			i += 2
			result += "\x1b[1m"

			for i+1 < len(content) && !(content[i] == '*' && content[i+1] == '*') {
				result += string(content[i])
				i++
			}
			result += "\x1b[0m"
			i += 2

		} else {
			result += string(content[i])
			i += 1
		}

	}
	return result
}

type ResponseMessage struct {
	Action  string `json:"action"`
	Message string `json: "message"`
	Code    string `json:"code"`
}

func parseJSON(content string) string {
	var v ResponseMessage
	json.Unmarshal([]byte(content), &v)
	return v.Message
}
