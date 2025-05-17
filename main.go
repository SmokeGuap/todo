package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoList []Task

func (t *TodoList) Add(title string) {
	task := Task{
		ID:        t.NextID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	*t = append(*t, task)
}

func (t *TodoList) NextID() int {
	maxID := 0

	for _, task := range *t {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1
}

func (t *TodoList) SetDone(index int, b bool) error {
	id := findIndexByID(*t, index)
	if id == -1 {
		return errors.New("invalid task ID")
	}
	(*t)[id].Done = b
	return nil
}

func (t *TodoList) Delete(index int) error {
	id := findIndexByID(*t, index)
	if id == -1 {
		return errors.New("invalid task ID")
	}
	*t = append((*t)[:id], (*t)[id+1:]...)
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
	id := findIndexByID(*t, index)
	if id == -1 {
		return errors.New("invalid task ID")
	}
	(*t)[id].Title = title
	return nil
}

func (t *TodoList) Print() {
	if len(*t) == 0 {
		fmt.Println("Your todo list is empty")
		return
	}

	todoPrint(t)
}

func (t *TodoList) Search(str string) (TodoList, error) {
	var matchedTodo TodoList
	for _, todo := range *t {
		if strings.Contains(strings.ToLower(todo.Title), strings.ToLower(str)) {
			matchedTodo = append(matchedTodo, todo)
		}
	}
	if len(matchedTodo) == 0 {
		return nil, fmt.Errorf("Tasks not found with title: %s", str)
	}

	return matchedTodo, nil
}

func (t *TodoList) Clear() {
	*t = nil
}

func (t *TodoList) Sort(key string) error {
	switch key {
	case "1":
		slices.SortFunc(*t, func(a, b Task) int {
			return cmp.Compare(a.Title, b.Title)
		})
	case "2":
		slices.SortFunc(*t, func(a, b Task) int {
			if a.CreatedAt.Before(b.CreatedAt) {
				return -1
			}
			if a.CreatedAt.After(b.CreatedAt) {
				return 1
			}
			return 0
		})
	default:
		return errors.New("Error: wrong sort field")

	}
	return nil
}

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
