# todo

A console task manager written in Go. It supports two modes: an interactive menu and command-line flags. Tasks are stored in a JSON file.

## Features

- View the task list with completion status and creation date
- Add, rename, and delete tasks
- Mark tasks as completed / uncompleted
- Search by title (case-insensitive)
- Sort by title or creation date
- Clear the entire list
- Automatic ID assignment (highest existing ID + 1)
- Data stored in a JSON file with a configurable name

## Build and run

```bash
# build
go build  .

# run
./todo
```

Or run directly without building:

```bash
go run .
```

## Usage

### Choosing the data file

On every launch, the program asks for the task file name:

```
Enter filename (default: todos.json):
```

- Empty input — `todos.json` is used
- If no extension is given, `.json` is appended automatically (e.g. `work` → `work.json`)
- If the file doesn't exist, it is created on the first save

### Interactive mode

Enabled by default (when no flags are passed) or with the `-interactive` flag. Available menu items:

| Item | Action                           |
| ---- | -------------------------------- |
| 1    | List all tasks                   |
| 2    | Add a task                       |
| 3    | Mark a task as completed         |
| 4    | Mark a task as uncompleted       |
| 5    | Change a task title              |
| 6    | Delete a task                    |
| 7    | Search by title                  |
| 8    | Clear the entire list            |
| 9    | Sort (by title or creation date) |
| 10   | Exit                             |

### Command-line mode

If at least one flag is passed (and `-interactive` is not set), the program performs the action and exits.

| Flag           | Description                |
| -------------- | -------------------------- |
| `-list`        | Print all tasks            |
| `-add "text"`  | Add a new task             |
| `-done ID`     | Mark a task as completed   |
| `-undone ID`   | Mark a task as uncompleted |
| `-delete ID`   | Delete a task              |
| `-interactive` | Start interactive mode     |

Examples:

```bash
# add a task
./todo -add "Buy something"

# show the list
./todo -list

# mark task with ID 2 as completed
./todo -done 2

# delete task with ID 3
./todo -delete 3

# add a task and print the list right away
./todo -add "Call friend" -list
```

> **Note:** `-done` and `-undone` cannot be used together — the program returns an error. Actions run in this order: `add` → `done`/`undone` → `delete` → `list`.
