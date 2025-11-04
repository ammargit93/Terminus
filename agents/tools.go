package agents

import (
	"fmt"
	"os"
)

func WriteFile(filename, content string) {
	docstring := `this function takes in filename and content and writes 
	the content to the file, creates the file
	if not created.
	`
	fmt.Println(docstring)
	os.WriteFile(filename, []byte(content), 0755)
}
