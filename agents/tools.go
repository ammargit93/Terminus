package agents

import (
	"os"
)

func WriteFile(args []string) string {
	os.WriteFile(args[0], []byte(args[1]), 0755)
	return ""
}
func MakeDirs(args []string) string {
	cwd, _ := os.Getwd()
	os.MkdirAll(cwd+"\\"+args[0], 0755)
	return ""
}
