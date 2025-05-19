package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

func getDataFilename(scanner *bufio.Scanner) string {
	fmt.Print("Enter filename (default: todos.json): ")

	var filename string

	if scanner.Scan() {
		filename = strings.TrimSpace(scanner.Text())
	}

	if filename == "" {
		filename = "todos.json"
	}

	if filepath.Ext(filename) == "" {
		filename += ".json"
	}

	return filename
}
