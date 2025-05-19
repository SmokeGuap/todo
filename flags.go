package main

import "flag"

type CLIFlags struct {
	List        bool
	Add         string
	Done        int
	Undone      int
	Delete      int
	Interactive bool
}

func parseFlags() CLIFlags {
	var flags CLIFlags

	flag.BoolVar(&flags.List, "list", false, "List all tasks")
	flag.StringVar(&flags.Add, "add", "", "Add new task")
	flag.IntVar(&flags.Done, "done", 0, "Mark task as completed")
	flag.IntVar(&flags.Undone, "undone", 0, "Mark task as uncompleted")
	flag.IntVar(&flags.Delete, "delete", 0, "Delete a task")
	flag.BoolVar(&flags.Interactive, "interactive", false, "Run interactive mode")

	flag.Parse()

	return flags
}
