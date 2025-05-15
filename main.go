package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoList []Task

func (t *TodoList) Add(title string) {
	task := Task{
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	*t = append(*t, task)
}

func (t *TodoList) SetDone(index int, b bool) error {
	if index < 0 || index >= len(*t) {
		return errors.New("invalid task index")
	}
	(*t)[index].Done = b
	return nil
}

func (t *TodoList) Delete(index int) error {
	if index < 0 || index >= len(*t) {
		return errors.New("invalid task index")
	}
	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}

func (t *TodoList) Save(filename string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *TodoList) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, t)
}

func (t *TodoList) Change(index int, title string) error {
	if index < 0 || index >= len(*t) {
		return errors.New("invalid task index")
	}
	(*t)[index].Title = title
	return nil
}

func (t *TodoList) Print() {
	if len(*t) == 0 {
		fmt.Println("Your todo list is empty")
		return
	}

	fmt.Println("Task list:")
	for i, task := range *t {
		status := "[ ]"
		if task.Done {
			status = "[\u2713]"
		}
		fmt.Printf("%d. %s %s (created: %v)\n", i+1, status, task.Title, task.CreatedAt.Format("2006.01.02 15:04"))
	}
}

func waitForEnter(s *bufio.Scanner) {
	fmt.Print("\nPress Enter to continue...")
	s.Scan()
}

func readTodoId(s *bufio.Scanner) (int, error) {
	fmt.Print("Write id: ")
	if !s.Scan() {
		return 0, errors.New("failed to read input")
	}
	id := strings.TrimSpace(s.Text())
	num, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("please enter a valid number")
	}
	return num - 1, nil
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

	err = todos.SetDone(id, true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		waitForEnter(s)
		return
	}

	if err := todos.Save(filename); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	} else if done {
		fmt.Println("Task marked as completed")
	} else {
		fmt.Println("Task marked as uncompleted")
	}

	waitForEnter(s)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter filename (default: todos.json): ")
	var dataFile string
	if scanner.Scan() {
		dataFile = strings.TrimSpace(scanner.Text())
	}
	if dataFile == "" {
		dataFile = "todos.json"
	}

	if filepath.Ext(dataFile) == "" {
		dataFile = dataFile + ".json"
	}

	fmt.Printf("File '%s' will be used.\n", dataFile)

	todos := &TodoList{}
	if err := todos.Load(dataFile); err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		os.Exit(1)
	}

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
		fmt.Println("7. Exit")
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

			err = todos.Delete(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				waitForEnter(scanner)
				continue
			}
			if err := todos.Save(dataFile); err != nil {
				fmt.Printf("Error saving data: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully")
			}
			waitForEnter(scanner)

		case "7":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Please choose a number from 1 to 7.")
			waitForEnter(scanner)
		}

	}
}
