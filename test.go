package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("COHERE_API_KEY"))
}
