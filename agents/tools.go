package agents

import (
	"os"
)

func WriteFile(args []string) string {
	finalPath := getCWD(args[0])
	os.WriteFile(finalPath, []byte(args[1]), 0755)
	return ""
}

func MakeDirs(args []string) string {
	finalPath := getCWD(args[0])
	os.MkdirAll(finalPath, 0755)
	return ""
}
