package model

import (
	"strings"
	"time"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	Priority    Priority  `json:"priority"`
	Project     string    `json:"project"` // e.g., "atlas"
	Contexts    []string  `json:"contexts"` // e.g., "@home", "@work"
	Category    string    `json:"category"`
}

func NewTask(title string) Task {
	return Task{
		Title:     title,
		CreatedAt: time.Now(),
		Priority:  PriorityMedium,
		Done:      false,
	}
}

func ParseTask(input string) Task {
	t := NewTask(input)

	// 1. Priority (Single)
	if strings.Contains(input, "!high") {
		t.Priority = PriorityHigh
		input = strings.ReplaceAll(input, "!high", "")
	} else if strings.Contains(input, "!med") {
		t.Priority = PriorityMedium
		input = strings.ReplaceAll(input, "!med", "")
	} else if strings.Contains(input, "!low") {
		t.Priority = PriorityLow
		input = strings.ReplaceAll(input, "!low", "")
	}

	// 2. Category (Single - First one wins, others removed to keep title clean)
	words := strings.Fields(input)
	cleanWords := []string{}
	foundCategory := false
	
	for _, w := range words {
		if strings.HasPrefix(w, "@") && len(w) > 1 {
			if !foundCategory {
				t.Category = strings.TrimPrefix(w, "@")
				foundCategory = true
			}
			// We skip appending this word to cleanWords, effectively removing it from title
		} else {
			cleanWords = append(cleanWords, w)
		}
	}

	t.Title = strings.Join(cleanWords, " ")
	return t
}
