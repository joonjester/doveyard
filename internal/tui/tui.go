package tui

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type email struct {
	sender     string
	receiver   string
	subject    string
	content    string
	time       time.Time
	readStatus bool
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
			sender:     "test@test.de",
			receiver:   "test2@test.de",
			subject:    "This is a test email",
			content:    "Hello World!",
			time:       time.Now(),
			readStatus: false,
		}, {
			sender:     "test@test.de",
			receiver:   "test2@test.de",
			subject:    "This is a second test email",
			content:    "Hello World!",
			time:       time.Now(),
			readStatus: false,
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
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	normalStyle := lipgloss.NewStyle()
	unreadStyle := lipgloss.NewStyle().Foreground(lipgloss.Cyan)
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.White)

	var lines string
	for i, email := range m.emails {
		line := fmt.Sprintf("[%s] %s\n %s", email.time.Format("02.01.2006"), email.sender, email.subject)

		_, ok := m.selected[i]
		if ok && email.readStatus == false {
			email.readStatus = true
		}

		styledLines := normalStyle.Render(line)
		if email.readStatus == false {
			styledLines = unreadStyle.Render(line)
		}

		if m.cursor == i {
			lines += selectedStyle.Render(styledLines) + "\n"
		} else {
			lines += styledLines + "\n"
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
