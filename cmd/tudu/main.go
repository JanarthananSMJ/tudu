package main

import (
	"fmt"
	"os"
	"path/filepath"

	"tudu/internal/commands"
	"tudu/internal/storage"
	"tudu/internal/tui"
)

func main() {
	storePath, err := defaultTodosPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve storage path: %v\n", err)
		os.Exit(1)
	}

	repo := storage.NewJSONRepository(storePath)
	svc := commands.NewService(repo)
	app := tui.NewModel(svc)

	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func defaultTodosPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tudu", "todos.json"), nil
}
