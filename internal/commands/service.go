package commands

import (
	"fmt"
	"strings"
	"time"

	"tudu/internal/models"
	"tudu/internal/storage"
)

// Service coordinates todo use-cases.
type Service struct {
	repo storage.Repository
}

func NewService(repo storage.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]models.Todo, error) {
	return s.repo.Load()
}

func (s *Service) Add(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	todos, err := s.repo.Load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	todos = append(todos, models.Todo{
		ID:        fmt.Sprintf("%d", now.UnixNano()),
		Title:     title,
		Completed: false,
		CreatedAt: now,
	})

	return s.repo.Save(todos)
}

// SaveAll overwrites persisted todos with the provided list.
func (s *Service) SaveAll(todos []models.Todo) error {
	return s.repo.Save(todos)
}
