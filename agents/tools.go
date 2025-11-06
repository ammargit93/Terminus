package agents

import (
	"os"
)

func WriteFile(args []string) string {
	os.WriteFile(args[0], []byte(args[1]), 0755)
	return ""
}
