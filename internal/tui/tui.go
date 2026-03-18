package tui

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type email struct {
	sender   string
	receiver string
	subject  string
	content  string
	time     time.Time
}

type size struct {
	width  int
	height int
}

type model struct {
	emails   []email
	cursor   int
	selected map[int]struct{}
	size     size
}

func initialModal() model {
	return model{
		emails: []email{{
			sender:   "test@test.de",
			receiver: "test2@test.de",
			subject:  "This is a test email",
			content:  "Hello World!",
			time:     time.Now(),
		}, {
			sender:   "test@test.de",
			receiver: "test2@test.de",
			subject:  "This is a second test email",
			content:  "Hello World!",
			time:     time.Now(),
		}},

		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.size.width = msg.Width
		m.size.height = msg.Height

	case tea.KeyPressMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.emails)-1 {
				m.cursor++
			}

		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	container := lipgloss.NewStyle().
		Width(m.size.width).
		Height(m.size.height).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	normalStyle := lipgloss.NewStyle()
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))

	lines := "What should we buy at the market?\n\n"

	for i, email := range m.emails {
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		line := fmt.Sprintf("[%s] %s\n %s", checked, email.subject, email.content)

		if m.cursor == i {
			lines += selectedStyle.Render(line) + "\n"
		} else {
			lines += normalStyle.Render(line) + "\n"
		}
	}

	lines += "\nPress q to quit.\n"

	return tea.NewView(container.Render(lines))
}

func Run() {
	p := tea.NewProgram(initialModal())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
