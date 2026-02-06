package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"atlas.todo/internal/model"
	"atlas.todo/internal/storage"
)

type state int

const (
	browsing state = iota
	adding
	searching
	deleting
	editing
)

type Grouping int

const (
	GroupNone Grouping = iota
	GroupCategory
	GroupDay
	GroupPriority
)

type Model struct {
	store        *storage.Store
	cursor       int
	state        state
	textInput    textinput.Model
	searchInput  textinput.Model
	sortByDate   bool
	sortAsc      bool
	showDone     bool
	grouping     Grouping
	width        int
	height       int
	taskToDelete model.Task
	taskToEdit   model.Task
	err          error
}

func NewModel(store *storage.Store) Model {
	ti := textinput.New()
	ti.Placeholder = "New task... (e.g. Buy milk @store @urgent !high)"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	si := textinput.New()
	si.Placeholder = "Search tasks..."
	si.CharLimit = 50
	si.Width = 30

	return Model{
		store:       store,
		textInput:   ti,
		searchInput: si,
		state:       browsing,
		sortByDate:  false,
		sortAsc:     false, // Newest first by default
		showDone:    true,
		grouping:    GroupNone,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		appStyle = appStyle.Width(msg.Width - 4).Height(msg.Height - 2)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case browsing:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				max := len(m.filteredTasks()) - 1
				if m.cursor < max {
					m.cursor++
				}
			case " ":
				tasks := m.filteredTasks()
				if len(tasks) > 0 && m.cursor < len(tasks) {
					targetID := tasks[m.cursor].ID
					for i, t := range m.store.Tasks {
						if t.ID == targetID {
							m.store.Toggle(i)
							break
						}
					}
					_ = m.store.Save()
				}
			case "d":
				tasks := m.filteredTasks()
				if len(tasks) > 0 && m.cursor < len(tasks) {
					m.taskToDelete = tasks[m.cursor]
					m.state = deleting
				}
			case "n":
				m.state = adding
				m.textInput.Reset()
				m.textInput.Focus()
				return m, textinput.Blink
			case "e":
				tasks := m.filteredTasks()
				if len(tasks) > 0 && m.cursor < len(tasks) {
					m.taskToEdit = tasks[m.cursor]
					m.state = editing
					m.textInput.SetValue(m.taskToEdit.Format())
					m.textInput.Focus()
					return m, textinput.Blink
				}
			case "/":
				m.state = searching
				m.searchInput.Reset()
				m.searchInput.Focus()
				return m, textinput.Blink
			case "s":
				if !m.sortByDate {
					m.sortByDate = true
					m.sortAsc = true // Default to Asc after first press
				} else if m.sortAsc {
					m.sortAsc = false // Switch to Desc
				} else {
					m.sortByDate = false // Back to Default
				}
				return m, nil
			case "c":
				m.showDone = !m.showDone
				m.cursor = 0
				return m, nil
			case "g":
				m.grouping++
				if m.grouping > GroupPriority {
					m.grouping = GroupNone
				}
				m.cursor = 0
				return m, nil
			}

		case adding:
			switch msg.String() {
			case "enter":
				text := m.textInput.Value()
				if text != "" {
					task := model.ParseTask(text)
					m.store.Add(task)
					_ = m.store.Save()
				}
				m.state = browsing
				m.cursor = 0
				return m, nil
			case "esc":
				m.state = browsing
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd

		case searching:
			switch msg.String() {
			case "enter", "esc":
				m.state = browsing
				m.cursor = 0
				return m, nil
			}
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd

		case deleting:
			switch msg.String() {
			case "y", "Y", "enter":
				targetID := m.taskToDelete.ID
				for i, t := range m.store.Tasks {
					if t.ID == targetID {
						m.store.Delete(i)
						break
					}
				}
				_ = m.store.Save()
				m.state = browsing
				tasks := m.filteredTasks()
				if m.cursor >= len(tasks) && m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "n", "N", "esc", "q":
				m.state = browsing
				return m, nil
			}

		case editing:
			switch msg.String() {
			case "enter":
				text := m.textInput.Value()
				if text != "" {
					updatedTask := model.ParseTask(text)
					for i, t := range m.store.Tasks {
						if t.ID == m.taskToEdit.ID {
							m.store.Tasks[i].Title = updatedTask.Title
							m.store.Tasks[i].Category = updatedTask.Category
							m.store.Tasks[i].Priority = updatedTask.Priority
							break
						}
					}
					_ = m.store.Save()
				}
				m.state = browsing
				return m, nil
			case "esc":
				m.state = browsing
				return m, nil
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) filteredTasks() []model.Task {
	var filtered []model.Task
	query := strings.ToLower(m.searchInput.Value())

	for _, t := range m.store.Tasks {
		// Filter by 'showDone'
		if !m.showDone && t.Done {
			continue
		}
		// Filter by search query
		if query == "" || strings.Contains(strings.ToLower(t.Title), query) {
			filtered = append(filtered, t)
		}
	}

	// 1. Grouping Sorts
	switch m.grouping {
	case GroupCategory:
		sort.SliceStable(filtered, func(i, j int) bool {
			cat1 := filtered[i].Category
			if cat1 == "" {
				cat1 = "Uncategorized"
			}
			cat2 := filtered[j].Category
			if cat2 == "" {
				cat2 = "Uncategorized"
			}
			
			if m.sortAsc {
				return cat1 < cat2
			}
			return cat1 > cat2
		})
	case GroupDay:
		sort.SliceStable(filtered, func(i, j int) bool {
			d1 := filtered[i].CreatedAt.Format("2006-01-02")
			d2 := filtered[j].CreatedAt.Format("2006-01-02")
			if m.sortAsc {
				return d1 < d2
			}
			return d1 > d2
		})
	case GroupPriority:
		sort.SliceStable(filtered, func(i, j int) bool {
			if m.sortAsc {
				return filtered[i].Priority < filtered[j].Priority
			}
			return filtered[i].Priority > filtered[j].Priority
		})
	}

	// 2. Date Sort
	if m.sortByDate && m.grouping != GroupDay {
		sort.SliceStable(filtered, func(i, j int) bool {
			if m.sortAsc {
				return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
			}
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		})
	}

	return filtered
}

func (m Model) View() string {
	if m.state == adding {
		return appStyle.Render(fmt.Sprintf(
			"Create a new task:\n\n%s\n\n(esc to cancel, enter to save)",
			m.textInput.View(),
		))
	}

	if m.state == editing {
		return appStyle.Render(fmt.Sprintf(
			"Edit task:\n\n%s\n\n(esc to cancel, enter to save)",
			m.textInput.View(),
		))
	}

	statusParts := []string{}
	orderStr := ""
	if m.sortByDate {
		if m.sortAsc {
			orderStr = "↑"
			statusParts = append(statusParts, "Sort: Asc "+orderStr)
		} else {
			orderStr = "↓"
			statusParts = append(statusParts, "Sort: Desc "+orderStr)
		}
	}

	if !m.showDone {
		statusParts = append(statusParts, "Hidden: Done")
	}

	switch m.grouping {
	case GroupCategory:
		statusParts = append(statusParts, "Group: Category")
	case GroupDay:
		statusParts = append(statusParts, "Group: Day")
	case GroupPriority:
		statusParts = append(statusParts, "Group: Priority")
	}
	
	statusStr := ""
	if len(statusParts) > 0 {
		statusStr = fmt.Sprintf(" [%s]", strings.Join(statusParts, ", "))
	}

	searchBar := ""
	if m.state == searching || m.searchInput.Value() != "" {
		searchBar = "\nSearch: " + m.searchInput.View() + "\n"
	}

	header := titleStyle.Render("Atlas Todo") + statusStr + searchBar + "\n"

	if m.state == deleting {
		prompt := fmt.Sprintf("Delete \"%s\"? (y/n)", m.taskToDelete.Title)
		return appStyle.Render(header + "\n" + deleteWarnStyle.Render(prompt))
	}

	displayTasks := m.filteredTasks()
	s := ""
	var lastGroupKey string
	
	for i, task := range displayTasks {
		currentGroupKey := ""
		switch m.grouping {
		case GroupCategory:
			currentGroupKey = task.Category
			if currentGroupKey == "" {
				currentGroupKey = "Uncategorized"
			}
		case GroupDay:
			currentGroupKey = task.CreatedAt.Format("Monday, 02 Jan 2006")
		case GroupPriority:
			switch task.Priority {
			case model.PriorityHigh:
				currentGroupKey = "!!! High Priority"
			case model.PriorityMedium:
				currentGroupKey = "!!  Medium Priority"
			case model.PriorityLow:
				currentGroupKey = "!   Low Priority"
			}
		}

		if m.grouping != GroupNone && currentGroupKey != lastGroupKey {
			s += groupHeaderStyle.Render(currentGroupKey) + "\n"
			lastGroupKey = currentGroupKey
		}

		// Selection & Base Style
		baseStyle := itemStyle
		if m.cursor == i {
			baseStyle = selectedItemStyle
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if task.Done {
			checked = "[x]"
		}

		// Compose the content without styles first
		catStr := ""
		if task.Category != "" {
			catStr = fmt.Sprintf(" (@%s)", task.Category)
		}
		dateStr := task.CreatedAt.Format(" (2006-01-02 15:04)")
		
		// Build the line with targeted styling
		var titlePart, catPart, datePart string
		
		if task.Done {
			// If done, everything is grey/strikethrough
			titlePart = doneStyle.Render(task.Title)
			catPart = doneStyle.Render(catStr)
			datePart = doneStyle.Render(dateStr)
		} else {
			// Normal state: use specific colors
			titlePart = task.Title
			catPart = categoryStyle.Render(catStr)
			datePart = dateStyle.Render(datePart) // Wait, was datePart or dateStr? Checking... dateStr.
			datePart = dateStyle.Render(dateStr)
		}

		// Assemble
		content := fmt.Sprintf("%s %s %s%s%s", cursor, checked, titlePart, catPart, datePart)
		
		// Apply outer selection style (padding/bold) but avoid overriding Foreground if already set
		s += baseStyle.Render(content) + "\n"
	}

	if len(displayTasks) == 0 {
		s = "\n  No tasks found.\n"
	}

		storageInfo := dateStyle.Render("\nStorage: ~/.atlas/todo.json (Auto-created)")

		help := helpStyle.Render("\nj/k: move • space: toggle • n: new • e: edit • /: search • d: delete • g: group • s: sort cycle • c: toggle done • q: quit")

	

		return appStyle.Render(header + s + storageInfo + help)

	}

	