package main

import (
	"fmt"
	"os"
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
	fmt.Println("  atlas.todo help          Show this help information")
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