package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	flags := parseFlags()

	scanner := bufio.NewScanner(os.Stdin)

	dataFile := getDataFilename(scanner)

	fmt.Printf("File '%s' will be used.\n", dataFile)

	todos := &TodoList{}
	if err := todos.Load(dataFile); err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		os.Exit(1)
	}

	if flag.NFlag() > 0 && !flags.Interactive {
		if err := runCLI(flags, todos, dataFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		return
	}

	runInteractive(scanner, todos, dataFile)
}
