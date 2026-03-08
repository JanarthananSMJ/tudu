package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tudu/internal/commands"
	"tudu/internal/models"
)

type todosLoadedMsg struct {
	todos []models.Todo
	err   error
}

type todosSavedMsg struct {
	err error
}

// Model is the Bubble Tea state container.
type Model struct {
	service    *commands.Service
	styles     uiStyles
	keyMap     KeyMap
	todos      []models.Todo
	cursor     int
	loading    bool
	adding     bool
	editing    bool
	editingIdx int
	textInput  textinput.Model
	err        error
	status     string
}

func NewModel(service *commands.Service) *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter todo title"
	ti.CharLimit = 200
	ti.Width = 50

	return &Model{
		service:    service,
		styles:     newUIStyles(),
		keyMap:     DefaultKeyMap(),
		todos:      []models.Todo{},
		cursor:     0,
		loading:    true,
		adding:     false,
		editing:    false,
		editingIdx: -1,
		textInput:  ti,
	}
}

func (m *Model) Start() error {
	_, err := tea.NewProgram(m).Run()
	return err
}

// Init loads todos from storage.
func (m *Model) Init() tea.Cmd {
	return m.loadTodosCmd()
}

// Update handles messages and key events.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case todosLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.status = "failed to load todos"
			return m, nil
		}
		m.todos = msg.todos
		if len(m.todos) == 0 {
			m.cursor = 0
		} else {
			m.cursor = clamp(m.cursor, 0, len(m.todos)-1)
		}
		m.status = fmt.Sprintf("loaded %d todo(s)", len(m.todos))
		return m, nil
	case todosSavedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "failed to save todos"
		}
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if !m.adding && !m.editing && msg.String() == m.keyMap.Quit {
			return m, tea.Quit
		}

		if m.loading {
			return m, nil
		}

		if m.adding {
			return m.handleAddInput(msg)
		}

		if m.editing {
			return m.handleEditInput(msg)
		}

		switch msg.String() {
		case m.keyMap.Down:
			if len(m.todos) > 0 && m.cursor < len(m.todos)-1 {
				m.cursor++
			}
			return m, nil
		case m.keyMap.Up:
			if len(m.todos) > 0 && m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case m.keyMap.Add:
			m.adding = true
			m.textInput.SetValue("")
			m.textInput.Focus()
			m.status = "add todo: type and press enter"
			return m, textinput.Blink
		case m.keyMap.Edit:
			if len(m.todos) == 0 {
				return m, nil
			}
			m.editing = true
			m.editingIdx = m.cursor
			m.textInput.SetValue(m.todos[m.cursor].Title)
			m.textInput.Focus()
			m.status = "edit todo: update title and press enter"
			return m, textinput.Blink
		case m.keyMap.Delete:
			if len(m.todos) == 0 {
				return m, nil
			}
			m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
			if len(m.todos) == 0 {
				m.cursor = 0
			} else if m.cursor >= len(m.todos) {
				m.cursor = len(m.todos) - 1
			}
			m.status = "todo deleted"
			return m, m.saveTodosCmd(cloneTodos(m.todos))
		case m.keyMap.Complete:
			if len(m.todos) == 0 {
				return m, nil
			}
			m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
			m.status = "todo completion toggled"
			return m, m.saveTodosCmd(cloneTodos(m.todos))
		default:
			return m, nil
		}
	}

	return m, nil
}

// View renders the TUI content.
func (m *Model) View() string {
	if m.loading {
		return m.styles.app.Render(m.styles.hint.Render("Loading todos..."))
	}

	lines := []string{
		m.styles.header.Render("TUI TODO"),
		m.styles.divider.Render(dividerLine(48)),
		"",
	}

	if len(m.todos) == 0 {
		lines = append(lines, m.styles.hint.Render("No todos yet. Press 'a' to add one."))
	} else {
		for i, todo := range m.todos {
			cursor := "  "
			if i == m.cursor {
				cursor = m.styles.cursor.Render("> ")
			}

			done := ""
			if todo.Completed {
				done = " ✓"
			}

			text := m.styles.todoText.Render(todo.Title + done)
			if todo.Completed {
				text = m.styles.todoDone.Render(todo.Title + done)
			}

			row := m.styles.row.Render(fmt.Sprintf("%s%s", cursor, text))
			if i == m.cursor {
				row = m.styles.rowSelected.Render(fmt.Sprintf("%-46s", row))
			}
			lines = append(lines, row)
		}
	}

	if m.editing {
		lines = append(lines, "", m.styles.hint.Render("Edit Todo"), m.textInput.View(), m.styles.hint.Render("enter: save | esc: cancel"))
	}
	if m.adding {
		lines = append(lines, "", m.styles.hint.Render("Add Todo"), m.textInput.View(), m.styles.hint.Render("enter: save | esc: cancel"))
	}

	if m.status != "" {
		lines = append(lines, "", m.styles.status.Render("Status: "+m.status))
	}
	if m.err != nil {
		lines = append(lines, m.styles.error.Render("Error: "+m.err.Error()))
	}

	if !m.editing && !m.adding {
		lines = append(lines,
			"",
			m.styles.divider.Render(dividerLine(48)),
			m.styles.footer.Render("j/k: navigate | a: add | d: delete | c: complete | e: edit | q: quit"),
		)
	}

	return m.styles.app.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (m *Model) handleEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.editing = false
		m.editingIdx = -1
		m.textInput.Blur()
		m.textInput.SetValue("")
		m.status = "edit canceled"
		return m, nil
	case tea.KeyEnter:
		title := strings.TrimSpace(m.textInput.Value())
		if title == "" {
			m.status = "title cannot be empty"
			return m, nil
		}
		if m.editingIdx < 0 || m.editingIdx >= len(m.todos) {
			m.status = "no todo selected for edit"
			m.editing = false
			m.editingIdx = -1
			m.textInput.Blur()
			m.textInput.SetValue("")
			return m, nil
		}
		m.todos[m.editingIdx].Title = title
		m.status = "todo updated"

		m.editing = false
		m.editingIdx = -1
		m.textInput.Blur()
		m.textInput.SetValue("")
		return m, m.saveTodosCmd(cloneTodos(m.todos))
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m *Model) handleAddInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.adding = false
		m.textInput.Blur()
		m.textInput.SetValue("")
		m.status = "add canceled"
		return m, nil
	case tea.KeyEnter:
		title := strings.TrimSpace(m.textInput.Value())
		if title == "" {
			m.status = "title cannot be empty"
			return m, nil
		}

		now := time.Now().UTC()
		m.todos = append(m.todos, models.Todo{
			ID:        fmt.Sprintf("%d", now.UnixNano()),
			Title:     title,
			Completed: false,
			CreatedAt: now,
		})
		m.cursor = len(m.todos) - 1
		m.adding = false
		m.textInput.Blur()
		m.textInput.SetValue("")
		m.status = "todo added"
		return m, m.saveTodosCmd(cloneTodos(m.todos))
	default:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m *Model) loadTodosCmd() tea.Cmd {
	return func() tea.Msg {
		todos, err := m.service.List()
		return todosLoadedMsg{todos: todos, err: err}
	}
}

func (m *Model) saveTodosCmd(todos []models.Todo) tea.Cmd {
	return func() tea.Msg {
		err := m.service.SaveAll(todos)
		return todosSavedMsg{err: err}
	}
}

func cloneTodos(src []models.Todo) []models.Todo {
	out := make([]models.Todo, len(src))
	copy(out, src)
	return out
}

func clamp(v, min, max int) int {
	if max < min {
		return min
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
