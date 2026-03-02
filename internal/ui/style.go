package ui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6B50FF")).
			Padding(0, 1).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("212")).
				Bold(true)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	doneStyle = lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.Color("#6C6C6C"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	dateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	groupHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B3B3FF")).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(0).
			Underline(true)

	deleteWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D787")).
			Bold(true)
			
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)
	        
	checkboxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#585858"))
            
	checkedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D787"))
)
