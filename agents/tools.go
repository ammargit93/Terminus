package agents

import (
	"os"
)

func WriteFile(filename, content string) {
	os.WriteFile(filename, []byte(content), 0755)
}
