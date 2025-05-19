package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runInteractive(scanner *bufio.Scanner, todos *TodoList, dataFile string) {
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println("───────────────────────────────────────")
		fmt.Println("Menu:")
		fmt.Println("1. List all tasks")
		fmt.Println("2. Add a new task")
		fmt.Println("3. Mark task as completed")
		fmt.Println("4. Mark task as uncompleted")
		fmt.Println("5. Change title")
		fmt.Println("6. Delete a task")
		fmt.Println("7. Search")
		fmt.Println("8. Clear all tasks")
		fmt.Println("9. Sort")
		fmt.Println("10. Exit")
		fmt.Println("───────────────────────────────────────")

		fmt.Print("Choose an option: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Critical input error: %v\n", err)
				os.Exit(1)
			}
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			todos.Print()
			waitForEnter(scanner)

		case "2":
			fmt.Print("Write title: ")
			if !scanner.Scan() {
				return
			}
			title := strings.TrimSpace(scanner.Text())

			if title == "" {
				fmt.Println("Error: task title cannot be empty")
				waitForEnter(scanner)
				continue
			}

			todos.Add(title)
			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Task added successfully")
			}
			waitForEnter(scanner)

		case "3":
			handleSetDone(scanner, todos, dataFile, true)

		case "4":
			handleSetDone(scanner, todos, dataFile, false)

		case "5":
			if len(*todos) == 0 {
				fmt.Println("Your todo list is empty")
				waitForEnter(scanner)
				continue
			}
			todos.Print()

			id, err := readTodoId(scanner)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}

			fmt.Print("Write title: ")
			if !scanner.Scan() {
				continue
			}
			title := strings.TrimSpace(scanner.Text())

			if title == "" {
				fmt.Println("Error: title cannot be empty")
				waitForEnter(scanner)
				continue
			}

			err = todos.Change(id, title)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}
			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Task title changed successfully")
			}
			waitForEnter(scanner)
		case "6":
			if len(*todos) == 0 {
				fmt.Println("Your todo list is empty")
				waitForEnter(scanner)
				continue
			}
			todos.Print()

			id, err := readTodoId(scanner)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}

			fmt.Print("Are you sure to clear (y/n): ")
			if !scanner.Scan() {
				continue
			}

			yes := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if yes == "y" {
				err = todos.Delete(id)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					waitForEnter(scanner)
					continue
				}
			} else {
				continue
			}

			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully")
			}
			waitForEnter(scanner)

		case "7":
			fmt.Print("Write title: ")
			if !scanner.Scan() {
				continue
			}
			title := strings.TrimSpace(scanner.Text())

			if title == "" {
				fmt.Println("Error: title cannot be empty")
				waitForEnter(scanner)
				continue
			}

			matchedTodos, err := todos.Search(title)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}

			todoPrint(&matchedTodos)
			waitForEnter(scanner)

		case "8":
			fmt.Print("Are you sure to clear (y/n): ")
			if !scanner.Scan() {
				continue
			}

			yes := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if yes == "y" {
				todos.Clear()
			} else {
				continue
			}
			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Tasks deleted successfully")
			}
			waitForEnter(scanner)
		case "9":

			fmt.Println("Sort by:")
			fmt.Println("1. Title")
			fmt.Println("2. Created date")

			fmt.Print("Choose sort: ")
			if !scanner.Scan() {
				continue
			}
			key := strings.TrimSpace(scanner.Text())
			err := todos.Sort(key)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}
			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Tasks sorted successfully")
			}
			waitForEnter(scanner)
		case "10":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Please choose a number from 1 to 10.")
			waitForEnter(scanner)
		}

	}
}
