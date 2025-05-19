package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func findIndexByID(t TodoList, id int) int {
	for i, task := range t {
		if task.ID == id {
			return i
		}
	}

	return -1
}

func todoPrint(t *TodoList) {
	fmt.Println("Task list:")
	for _, task := range *t {
		status := "[ ]"
		if task.Done {
			status = "[\u2713]"
		}
		fmt.Printf(
			"%d. %s %s (created: %s)\n",
			task.ID,
			status,
			task.Title,
			task.CreatedAt.Format("2006.01.02 15:04"),
		)
	}
	fmt.Printf("Task count: %d\n", len(*t))
}

func waitForEnter(s *bufio.Scanner) {
	fmt.Print("\nPress Enter to continue...")
	s.Scan()
}

func readTodoId(s *bufio.Scanner) (int, error) {
	fmt.Print("Write ID: ")

	if !s.Scan() {
		return 0, errors.New("failed to read input")
	}

	id, err := strconv.Atoi(strings.TrimSpace(s.Text()))
	if err != nil {
		return 0, errors.New("please enter a valid number")
	}

	if id <= 0 {
		return 0, errors.New("ID must be positive")
	}

	return id, nil
}

func handleSetDone(s *bufio.Scanner,
	todos *TodoList,
	filename string,
	done bool) {
	if len(*todos) == 0 {
		fmt.Println("Your todo list is empty")
		waitForEnter(s)
		return
	}
	todos.Print()

	id, err := readTodoId(s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		waitForEnter(s)
		return
	}

	err = todos.SetDone(id, done)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		waitForEnter(s)
		return
	}

	if err := todos.Save(filename); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	} else if done {
		fmt.Printf("Task %d marked as completed\n", id)
	} else {
		fmt.Printf("Task %d marked as uncompleted\n", id)
	}

	waitForEnter(s)
}

func handleFlagSetDone(todos *TodoList, filename string, id int, done bool) error {
	if err := todos.SetDone(id, done); err != nil {
		return err
	}

	if err := todos.Save(filename); err != nil {
		return err
	}

	if done {
		fmt.Printf("Task %d marked as completed\n", id)
	} else {
		fmt.Printf("Task %d marked as uncompleted\n", id)
	}

	return nil
}
