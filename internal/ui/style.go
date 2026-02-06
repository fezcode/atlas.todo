package ui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

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
			Foreground(lipgloss.Color("#999999")).
			MarginTop(1)

	dateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999"))

	groupHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Bold(true).
			PaddingTop(1).
			PaddingBottom(0)

	deleteWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	categoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF"))
)
