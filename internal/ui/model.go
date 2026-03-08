package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
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
	showingHelp
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
	statusMsg    string
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
		sortByDate:  store.Config.SortByDate,
		sortAsc:     store.Config.SortAsc,
		showDone:    store.Config.ShowDone,
		grouping:    Grouping(store.Config.Grouping),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type clearStatusMsg struct{}

func clearStatus() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
				m.store.Config.SortByDate = m.sortByDate
				m.store.Config.SortAsc = m.sortAsc
				_ = m.store.Save()
				return m, nil
			case "y":
				tasks := m.filteredTasks()
				if len(tasks) > 0 && m.cursor < len(tasks) {
					task := tasks[m.cursor]
					content := task.Title
					if task.Category != "" {
						content = fmt.Sprintf("%s (@%s)", task.Title, task.Category)
					}
					_ = clipboard.WriteAll(content)
					m.statusMsg = "✓ Copied to clipboard!"
					return m, clearStatus()
				}
				return m, nil
			case "c":
				m.showDone = !m.showDone
				m.cursor = 0
				m.store.Config.ShowDone = m.showDone
				_ = m.store.Save()
				return m, nil
			case "g":
				m.grouping++
				if m.grouping > GroupPriority {
					m.grouping = GroupNone
				}
				m.cursor = 0
				m.store.Config.Grouping = int(m.grouping)
				_ = m.store.Save()
				return m, nil
			case "h":
				m.state = showingHelp
				return m, nil
			}

		case showingHelp:
			switch msg.String() {
			case "esc", "q", "h":
				m.state = browsing
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
	// 1. Padding Logic (Manual)
	topPad := 0
	if m.height > 15 { topPad = 1 }
	botPad := 0
	if m.height > 10 { botPad = 1 }

	style := appStyle.Width(m.width - 4)

	// 2. State-specific Views
	if m.state == adding {
		return style.Render(fmt.Sprintf(
			"Create a new task:\n\n%s\n\n(esc to cancel, enter to save)",
			m.textInput.View(),
		))
	}

	if m.state == editing {
		return style.Render(fmt.Sprintf(
			"Edit task:\n\n%s\n\n(esc to cancel, enter to save)",
			m.textInput.View(),
		))
	}

	if m.state == showingHelp {
		content := titleStyle.Render("Atlas Todo - Help & Tutorial") + "\n\n"
		
		content += groupHeaderStyle.Render("Metadata Basics") + "\n"
		content += "  • Category: Use @ (e.g., \"Buy milk @grocery\")\n"
		content += "  • Priority: Use ! (e.g., \"Fix bug !high\", \"!low\")\n"
		content += "  • Multiple: \"Meet John @work !medium\"\n\n"
		
		content += groupHeaderStyle.Render("Commands") + "\n"
		content += "  ↑/↓, j/k: move cursor • space: toggle done\n"
		content += "  n: new task      • e: edit selected\n"
		content += "  d: delete task   • y: copy to clipboard\n"
		content += "  s: cycle sort    • c: toggle completed\n"
		content += "  g: cycle groups  • /: search tasks\n"
		content += "  h: toggle help   • q: quit\n\n"

		content += helpStyle.Render("(press h or esc to return)")
		
		return style.PaddingTop(topPad).Render(content)
	}

	// 3. Header Construction
	statusParts := []string{}
	if m.sortByDate {
		orderStr := "↑"
		if !m.sortAsc { orderStr = "↓" }
		statusParts = append(statusParts, "Sort: "+orderStr)
	}
	if !m.showDone { statusParts = append(statusParts, "Hidden: Done") }
	switch m.grouping {
	case GroupCategory: statusParts = append(statusParts, "Group: Category")
	case GroupDay: statusParts = append(statusParts, "Group: Day")
	case GroupPriority: statusParts = append(statusParts, "Group: Priority")
	}
	
	statusStr := ""
	if len(statusParts) > 0 {
		statusStr = fmt.Sprintf(" [%s]", strings.Join(statusParts, ", "))
		if len(statusStr) > m.width-20 && m.width > 25 {
			statusStr = statusStr[:m.width-23] + "...]"
		} else if m.width <= 25 {
			statusStr = ""
		}
	}

	headerText := titleStyle.Render("Atlas Todo") + statusStr
	headerLines := 1
	
	searchBar := ""
	if m.state == searching || m.searchInput.Value() != "" {
		searchBar = "\nSearch: " + m.searchInput.View()
		headerLines++
	}
	headerText += searchBar + "\n" // Final newline
	headerLines++

	if m.state == deleting {
		prompt := fmt.Sprintf("Delete \"%s\"? (y/n)", m.taskToDelete.Title)
		return appStyle.Render(headerText + "\n" + deleteWarnStyle.Render(prompt))
	}

	// 4. Footer Logic & Height Budgeting
	footerHeight := 0
	showStatus := m.statusMsg != "" && m.height > 10
	showHelpLine := m.height > 6
	
	if showStatus { footerHeight += 1 }
	if showHelpLine { footerHeight += 1 }

	// availableHeight is what's left for the task list
	reserved := topPad + botPad + headerLines + footerHeight
	lineBudget := m.height - reserved
	if lineBudget < 1 { lineBudget = 1 }

	displayTasks := m.filteredTasks()
	if len(displayTasks) > 0 && m.cursor >= len(displayTasks) {
		m.cursor = len(displayTasks) - 1
	}

	// 5. Robust Scroll / Offset Calculation
	startIdx := m.cursor - (lineBudget / 2)
	if startIdx < 0 { startIdx = 0 }

	// Simulation loop to push startIdx forward until cursor is visible
	for {
		linesNeeded := 0
		if startIdx > 0 { linesNeeded++ } // hidden above
		
		tempLastGroup := ""
		if startIdx > 0 && m.grouping != GroupNone {
			prevTask := displayTasks[startIdx-1]
			switch m.grouping {
			case GroupCategory: tempLastGroup = prevTask.Category
			case GroupDay: tempLastGroup = prevTask.CreatedAt.Format("Monday, 02 Jan 2006")
			case GroupPriority:
				switch prevTask.Priority {
				case model.PriorityHigh: tempLastGroup = "!!! High Priority"
				case model.PriorityMedium: tempLastGroup = "!!  Medium Priority"
				case model.PriorityLow: tempLastGroup = "!   Low Priority"
				}
			}
			if tempLastGroup == "" && m.grouping == GroupCategory { tempLastGroup = "Uncategorized" }
		}

		cursorReached := false
		for i := startIdx; i < len(displayTasks); i++ {
			tBudget := lineBudget
			if i < len(displayTasks)-1 { tBudget-- } // reserve for hidden below

			if m.grouping != GroupNone {
				gKey := ""
				switch m.grouping {
				case GroupCategory: gKey = displayTasks[i].Category
					if gKey == "" { gKey = "Uncategorized" }
				case GroupDay: gKey = displayTasks[i].CreatedAt.Format("Monday, 02 Jan 2006")
				case GroupPriority:
					switch displayTasks[i].Priority {
					case model.PriorityHigh: gKey = "!!! High Priority"
					case model.PriorityMedium: gKey = "!!  Medium Priority"
					case model.PriorityLow: gKey = "!   Low Priority"
					}
				}
				if gKey != tempLastGroup {
					if linesNeeded > 0 { linesNeeded++ } // newline
					linesNeeded++ // header
					tempLastGroup = gKey
				}
			}
			linesNeeded++ // task
			
			if i == m.cursor {
				if linesNeeded <= tBudget { cursorReached = true }
				break
			}
			if linesNeeded >= tBudget { break }
		}

		if cursorReached || startIdx >= m.cursor { break }
		startIdx++
	}

	// 6. List Rendering (Budgeted)
	s := ""
	linesUsed := 0
	
	if startIdx > 0 && linesUsed < lineBudget {
		s += helpStyle.Render(fmt.Sprintf("  ... %d hidden above ...", startIdx)) + "\n"
		linesUsed++
	}

	var lastGroupKey string
	if startIdx > 0 && m.grouping != GroupNone {
		prevTask := displayTasks[startIdx-1]
		switch m.grouping {
		case GroupCategory: lastGroupKey = prevTask.Category
		case GroupDay: lastGroupKey = prevTask.CreatedAt.Format("Monday, 02 Jan 2006")
		case GroupPriority:
			switch prevTask.Priority {
			case model.PriorityHigh: lastGroupKey = "!!! High Priority"
			case model.PriorityMedium: lastGroupKey = "!!  Medium Priority"
			case model.PriorityLow: lastGroupKey = "!   Low Priority"
			}
		}
		if lastGroupKey == "" && m.grouping == GroupCategory { lastGroupKey = "Uncategorized" }
	}

	lastTaskIdx := startIdx - 1
	for i := startIdx; i < len(displayTasks); i++ {
		// Reserve 1 line for "hidden below" if not at the end
		currentBudget := lineBudget
		if i < len(displayTasks)-1 { currentBudget-- }

		if linesUsed >= currentBudget { break }

		task := displayTasks[i]
		
		// Handle Grouping
		if m.grouping != GroupNone {
			currentGroupKey := ""
			switch m.grouping {
			case GroupCategory:
				currentGroupKey = task.Category
				if currentGroupKey == "" { currentGroupKey = "Uncategorized" }
			case GroupDay:
				currentGroupKey = task.CreatedAt.Format("Monday, 02 Jan 2006")
			case GroupPriority:
				switch task.Priority {
				case model.PriorityHigh: currentGroupKey = "!!! High Priority"
				case model.PriorityMedium: currentGroupKey = "!!  Medium Priority"
				case model.PriorityLow: currentGroupKey = "!   Low Priority"
				}
			}
			
			if currentGroupKey != lastGroupKey {
				// No space for newline/header/task combo? 
				// We need at least 2 lines (header + task) or 3 (newline + header + task)
				neededLines := 2
				if linesUsed > 0 { neededLines = 3 }
				
				if linesUsed + neededLines <= currentBudget {
					if linesUsed > 0 {
						s += "\n"
						linesUsed++
					}
					s += groupHeaderStyle.Render(currentGroupKey) + "\n"
					lastGroupKey = currentGroupKey
					linesUsed++
				} else {
					break // No room for the next group
				}
			}
		}

		// Render Task
		if linesUsed < currentBudget {
			baseStyle := itemStyle
			if m.cursor == i { baseStyle = selectedItemStyle }

			cursor := " "
			if m.cursor == i { cursor = cursorStyle.Render("❯") }

			checked := checkboxStyle.Render("☐")
			if task.Done { checked = checkedStyle.Render("☑") }

			catStr := ""
			if task.Category != "" { catStr = fmt.Sprintf(" (@%s)", task.Category) }
			dateStr := task.CreatedAt.Format(" (2006-01-02 15:04)")
			
			var titlePart, catPart, datePart string
			if task.Done {
				titlePart = doneStyle.Render(task.Title)
				catPart = doneStyle.Render(catStr)
				datePart = doneStyle.Render(dateStr)
			} else {
				titlePart = task.Title
				catPart = categoryStyle.Render(catStr)
				datePart = dateStyle.Render(dateStr)
			}

			content := fmt.Sprintf("%s %s %s%s%s", cursor, checked, titlePart, catPart, datePart)
			s += baseStyle.Render(content) + "\n"
			linesUsed++
			lastTaskIdx = i
		} else {
			break
		}
	}

	if lastTaskIdx < len(displayTasks)-1 {
		hiddenBelow := len(displayTasks) - 1 - lastTaskIdx
		if hiddenBelow > 0 && linesUsed < lineBudget {
			s += helpStyle.Render(fmt.Sprintf("  ... %d hidden below ...", hiddenBelow)) + "\n"
			linesUsed++
		}
	}

	if len(displayTasks) == 0 {
		s = "\n  No tasks found.\n"
	}

	// 7. Footer Construction
	footer := ""
	if showStatus { footer += "\n" + statusStyle.Render(m.statusMsg) }
	
	if showHelpLine {
		footer += "\n" + helpStyle.Render("h: help")
	}

	// 8. Final Assembly
	res := strings.Repeat("\n", topPad) + headerText + s + footer
	return style.Render(res)
}

	