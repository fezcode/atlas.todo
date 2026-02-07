package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"atlas.todo/internal/model"
	"atlas.todo/internal/storage"
	"atlas.todo/internal/ui"
)

func main() {
	store, err := storage.NewStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing store: %v\n", err)
		os.Exit(1)
	}

	if err := store.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	// CLI Mode: Handle arguments
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "add":
			if len(os.Args) < 3 {
				fmt.Println("Usage: atlas.todo add \"Task text\"")
				os.Exit(1)
			}
			text := strings.Join(os.Args[2:], " ")
			task := model.ParseTask(text)
			
			store.Add(task)
			if err := store.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving task: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Task added: %s\n", task.Title)
			return
		case "list":
			// Default values
			limit := 5
			sortOrder := "default"

			// Parse arguments
			if len(os.Args) > 2 {
				for _, arg := range os.Args[2:] {
					if arg == "asc" || arg == "desc" {
						sortOrder = arg
					} else {
						if val, err := strconv.Atoi(arg); err == nil {
							limit = val
						}
					}
				}
			}

			// Filter pending tasks
			var pending []model.Task
			for _, t := range store.Tasks {
				if !t.Done {
					pending = append(pending, t)
				}
			}

			// Sort
			if sortOrder == "asc" {
				sort.Slice(pending, func(i, j int) bool {
					return pending[i].Priority < pending[j].Priority
				})
			} else if sortOrder == "desc" {
				sort.Slice(pending, func(i, j int) bool {
					return pending[i].Priority > pending[j].Priority
				})
			}

			// Print
			if len(pending) == 0 {
				fmt.Println("No pending tasks! 🎉")
				return
			}

			count := 0
			for _, t := range pending {
				if count >= limit {
					break
				}

				prioMarker := " "
				if t.Priority == model.PriorityHigh {
					prioMarker = "!"
				} else if t.Priority == model.PriorityLow {
					prioMarker = "."
				}

				fmt.Printf("[%s] %s\n", prioMarker, t.Format())
				count++
			}
			return
		case "help", "--help", "-h":
			showHelp()
			return
		}
	}

	// TUI Mode
	p := tea.NewProgram(ui.NewModel(store), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Atlas Todo - A fast, minimalist task manager for your terminal.")
	fmt.Println("\nUsage:")
	fmt.Println("  atlas.todo               Start the interactive TUI")
	fmt.Println("  atlas.todo add \"[task]\"  Quickly add a task via command line")
	fmt.Println("  atlas.todo list [opts]   List pending tasks (useful for MOTD)")
	fmt.Println("  atlas.todo help          Show this help information")
	fmt.Println("\nList Options:")
	fmt.Println("  atlas.todo list          Show top 5 tasks (default order)")
	fmt.Println("  atlas.todo list 3        Show top 3 tasks")
	fmt.Println("  atlas.todo list asc      Show tasks sorted by priority (Low -> High)")
	fmt.Println("  atlas.todo list desc 10  Show top 10 tasks sorted by priority (High -> Low)")
	fmt.Println("\nNote: When using 'add' from CLI, wrap your task in quotes if it contains")
	fmt.Println("      special characters or metadata like @category or !priority.")
	fmt.Println("      Example: atlas.todo add \"Buy milk @grocery !high\"")
	fmt.Println("\nTUI Controls:")
	fmt.Println("  j/k, up/down   Navigate through tasks")
	fmt.Println("  Space          Toggle task completion")
	fmt.Println("  n              Add a new task")
	fmt.Println("  /              Search tasks")
	fmt.Println("  d              Delete selected task")
	fmt.Println("  s              Toggle sort by date added")
	fmt.Println("  c              Toggle showing completed tasks")
	fmt.Println("  q, esc         Quit the application")
	fmt.Println("\nMetadata Parsing:")
	fmt.Println("  Include '!high' in your task title to set high priority.")
	fmt.Println("\nStorage:")
	fmt.Println("  Tasks are stored locally in ~/.atlas/todo.json.")
	fmt.Println("  The directory and file will be created automatically on first run.")
}