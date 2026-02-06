package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"atlas.todo/internal/model"
)

type Store struct {
	mu       sync.Mutex
	filePath string
	Tasks    []model.Task
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

	return json.Unmarshal(data, &s.Tasks)
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.Tasks, "", "  ")
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
