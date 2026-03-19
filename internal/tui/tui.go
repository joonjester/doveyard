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

type email struct {
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

type model struct {
	emails    map[int]email
	cursor    int
	selected  map[int]struct{}
	size      size
	showEmail bool
	viewType  ViewType
}

func initialModal() model {
	return model{
		emails: map[int]email{
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

		viewType: horizontalSplit,
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
			email := m.emails[m.cursor]
			_, ok := m.selected[m.cursor]
			if ok {
				email.hasRead = true
				delete(m.selected, m.cursor)
			} else {
				email.hasRead = false
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	width := m.size.width
	height := m.size.height

	fullWidth := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(width).
		Height(height)

	emailInbox, openedEmailIndex := allEmails(m)
	if openedEmailIndex != nil {
		email := m.emails[*openedEmailIndex]
		openedEmail := emailContent(email, width)

		display := splitView(m, emailInbox, openedEmail)

		return tea.NewView(display)
	}

	return tea.NewView(fullWidth.Render(emailInbox))
}

func getSplitStyle(isVertical bool, width int, height int) lipgloss.Style {
	if isVertical {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Width(width).
			Height(height / 2)
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(width / 2).
		Height(height)
}

func splitView(m model, emailInbox string, openEmail string) string {
	isVertical := m.viewType == verticalSplit
	viewStyleRender := getSplitStyle(isVertical, m.size.width, m.size.height)

	emailInboxRender := viewStyleRender.Render(emailInbox)
	emailContentRender := viewStyleRender.Render(openEmail)

	if isVertical {
		return lipgloss.JoinVertical(lipgloss.Left, emailInboxRender, emailContentRender)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, emailInboxRender, emailContentRender)
}

func emailContent(email email, width int) string {
	emailAddressStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	subjectStyle := lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.Border{Bottom: "="}).
		BorderBottom(true).Width((width / 2) - 2)

	var emailContent string

	addresses := fmt.Sprintf("From: %s\nTo: %s\n", email.sender, email.receiver)
	emailContent += fmt.Sprintf(
		"%s\n%s\n%s",
		emailAddressStyle.Render(addresses),
		subjectStyle.Render(email.subject),
		email.content,
	)

	return emailContent
}

func allEmails(m model) (string, *int) {
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
