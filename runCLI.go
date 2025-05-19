package main

import (
	"errors"
	"fmt"
)

func runCLI(flags CLIFlags, todos *TodoList, dataFile string) error {
	if flags.Add != "" {
		todos.Add(flags.Add)
		if err := todos.Save(dataFile); err != nil {
			return err
		} else {
			fmt.Println("Task added successfully")
		}
	}

	if flags.Done != 0 && flags.Undone != 0 {
		return errors.New("use either -done or -undone")
	}

	if flags.Done != 0 {
		if err := handleFlagSetDone(todos, dataFile, flags.Done, true); err != nil {
			return err
		}
	}

	if flags.Undone != 0 {
		if err := handleFlagSetDone(todos, dataFile, flags.Undone, false); err != nil {
			return err
		}
	}

	if flags.Delete != 0 {
		err := todos.Delete(flags.Delete)
		if err != nil {
			return err
		}
		if err := todos.Save(dataFile); err != nil {
			return err
		} else {
			fmt.Println("Task deleted successfully")
		}
	}

	if flags.List {
		todos.Print()
	}

	return nil
}
