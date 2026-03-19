package tui

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ViewType string

const (
	horizontalSplit ViewType = "horizontalSplit"
	verticalSplit   ViewType = "verticalSplit"
)

type Email struct {
	sender   string
	receiver string
	subject  string
	content  string
	time     time.Time
	hasRead  bool
}

type size struct {
	width  int
	height int
}

type Model struct {
	emails    map[int]Email
	cursor    int
	selected  map[int]struct{}
	size      size
	showEmail bool
	viewType  ViewType
}

func initialModal() Model {
	return Model{
		emails: map[int]Email{
			0: {
				sender:   "test@test.de",
				receiver: "test2@test.de",
				subject:  "This is a test email",
				content:  "Hi World!",
				time:     time.Now(),
				hasRead:  false,
			},
			1: {
				sender:   "test@test.de",
				receiver: "test2@test.de",
				subject:  "This is a second test email",
				content:  "Hello World!",
				time:     time.Now(),
				hasRead:  false,
			},
		},

		viewType: verticalSplit,
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			email := m.emails[m.cursor]
			_, ok := m.selected[m.cursor]
			if ok {
				email.hasRead = true
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	width := m.size.width
	height := m.size.height

	fullWidth := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(width).
		Height(height)

	emailInbox, openedEmailIndex := allEmails(m)
	if openedEmailIndex == nil {
		return tea.NewView(fullWidth.Render(emailInbox))
	}
	openedEmail := EmailContent(m.emails[*openedEmailIndex], width)
	display := SplitView(m, emailInbox, openedEmail)

	return tea.NewView(display)
}

func allEmails(m Model) (string, *int) {
	unreadEmail := lipgloss.NewStyle().Foreground(lipgloss.Cyan)
	hoveringEmail := lipgloss.NewStyle().Background(lipgloss.White)
	defaultEmail := lipgloss.NewStyle()

	var emails string
	var openedEmailIndex *int
	for i, email := range m.emails {
		message := fmt.Sprintf("[%s] %s\n %s",
			email.time.Format("02.01.2006"),
			email.sender,
			email.subject)

		_, ok := m.selected[i]
		if ok {
			openedEmailIndex = &i
		}

		if !email.hasRead && !ok {
			message = unreadEmail.Render(message)
		} else {
			message = defaultEmail.Render(message)
		}

		if m.cursor == i {
			message = hoveringEmail.Render(message) + "\n"
		} else {
			message = defaultEmail.Render(message) + "\n"
		}

		emails += message
	}

	return emails, openedEmailIndex
}

func Run() {
	p := tea.NewProgram(initialModal())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
