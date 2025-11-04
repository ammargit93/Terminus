package main

import (
	"encoding/json"
	"fmt"
)

func ParseJSON(content string) {
	var v map[string]any
	json.Unmarshal([]byte(content), &v)
	fmt.Println(v)
}

// func main() {
// 	content := `
// 	{
//         "message": "Writing 'hello world' to a text file",
//         "action": "write_file",
//         "code": "with open('hello.txt', 'w') as f: f.write('hello world')"
// 	}
// 	`
// 	ParseJSON(content)
// }
