package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"tudu/internal/models"
)

const (
	tudoDirName   = ".tudu"
	todosFileName = "todos.json"
)

// Repository defines persistence operations for todos.
type Repository interface {
	Load() ([]models.Todo, error)
	Save([]models.Todo) error
}

// JSONRepository stores todos in ~/.tudu/todos.json.
type JSONRepository struct {
	path string
}

type fileSchema struct {
	Todos []models.Todo `json:"todos"`
}

// NewJSONRepository creates a repository at ~/.tudu/todos.json when path is empty.
func NewJSONRepository(path string) *JSONRepository {
	if path == "" {
		path = defaultPath()
	}
	return &JSONRepository{path: path}
}

// Load returns all todos from disk.
func (r *JSONRepository) Load() ([]models.Todo, error) {
	return r.loadTodos()
}

// Save persists all todos to disk.
func (r *JSONRepository) Save(todos []models.Todo) error {
	return r.saveTodos(todos)
}

func (r *JSONRepository) loadTodos() ([]models.Todo, error) {
	if err := r.ensureFile(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []models.Todo{}, nil
	}

	var schema fileSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	if schema.Todos == nil {
		return []models.Todo{}, nil
	}

	return schema.Todos, nil
}

func (r *JSONRepository) saveTodos(todos []models.Todo) error {
	if err := r.ensureFile(); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(fileSchema{Todos: todos}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, payload, 0o644)
}

func (r *JSONRepository) ensureFile() error {
	if r.path == "" {
		return errors.New("storage path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(r.path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	initial, err := json.MarshalIndent(fileSchema{Todos: []models.Todo{}}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, initial, 0o644)
}

func defaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(tudoDirName, todosFileName)
	}
	return filepath.Join(home, tudoDirName, todosFileName)
}
