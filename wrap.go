package main

import (
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

		// Detect bullet or indent prefix (like "* ", "- ", "> ", etc.)
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

	return strings.Join(wrapped, "\n")
}
