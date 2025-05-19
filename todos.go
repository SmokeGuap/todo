package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
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

func (t *TodoList) SetDone(id int, done bool) error {
	index := findIndexByID(*t, id)
	if index == -1 {
		return errors.New("invalid task ID")
	}
	(*t)[index].Done = done
	return nil
}

func (t *TodoList) Delete(id int) error {
	index := findIndexByID(*t, id)
	if index == -1 {
		return errors.New("invalid task ID")
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

func (t *TodoList) Change(id int, title string) error {
	index := findIndexByID(*t, id)
	if index == -1 {
		return errors.New("invalid task ID")
	}
	(*t)[index].Title = title
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
	str = strings.ToLower(str)
	var matchedTodo TodoList
	for _, todo := range *t {
		if strings.Contains(strings.ToLower(todo.Title), str) {
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
