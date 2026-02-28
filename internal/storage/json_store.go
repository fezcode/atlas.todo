package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"atlas.todo/internal/model"
)

type Config struct {
	ShowDone   bool `json:"show_done"`
	SortByDate bool `json:"sort_by_date"`
	SortAsc    bool `json:"sort_asc"`
	Grouping   int  `json:"grouping"`
}

type storeData struct {
	Tasks  []model.Task `json:"tasks"`
	Config Config       `json:"config"`
}

type Store struct {
	mu       sync.Mutex
	filePath string
	Tasks    []model.Task
	Config   Config
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".atlas")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	return &Store{
		filePath: filepath.Join(configDir, "todo.json"),
		Tasks:    []model.Task{},
		Config: Config{
			ShowDone:   false,
			SortByDate: false,
			SortAsc:    false,
			Grouping:   0,
		},
	}, nil
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil // New store, no file yet
	}
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return nil
	}

	// Case 1: New format (Object)
	if trimmed[0] == '{' {
		var sd storeData
		if err := json.Unmarshal(data, &sd); err != nil {
			return fmt.Errorf("failed to parse as object: %w", err)
		}
		s.Tasks = sd.Tasks
		s.Config = sd.Config
		return nil
	}

	// Case 2: Old format (Array)
	if trimmed[0] == '[' {
		var tasks []model.Task
		if err := json.Unmarshal(data, &tasks); err != nil {
			return fmt.Errorf("failed to parse as array: %w", err)
		}
		s.Tasks = tasks
		return nil
	}

	return fmt.Errorf("unknown file format (must be { or [)")
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sd := storeData{
		Tasks:  s.Tasks,
		Config: s.Config,
	}

	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *Store) Add(t model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Simple ID gen if not present (UUID would be better but keeping deps low for now)
	if t.ID == "" {
		t.ID = time.Now().Format("20060102150405")
	}
	s.Tasks = append(s.Tasks, t)
}

func (s *Store) Toggle(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index >= 0 && index < len(s.Tasks) {
		s.Tasks[index].Done = !s.Tasks[index].Done
	}
}

func (s *Store) Delete(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index >= 0 && index < len(s.Tasks) {
		s.Tasks = append(s.Tasks[:index], s.Tasks[index+1:]...)
	}
}
