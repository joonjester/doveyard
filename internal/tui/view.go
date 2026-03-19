package tui

import (
	"charm.land/lipgloss/v2"
	"fmt"
)

func getSplitStyle(isVertical bool, width int, height int) lipgloss.Style {
	w, h := width/2, height
	if isVertical {
		w, h = width, height/2
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(w).
		Height(h)
}

func SplitView(m Model, emailInbox string, openEmail string) string {
	isVertical := m.viewType == verticalSplit
	viewStyleRender := getSplitStyle(isVertical, m.size.width, m.size.height)

	emailInbox, emailContent := viewStyleRender.Render(emailInbox), viewStyleRender.Render(openEmail)

	if isVertical {
		return lipgloss.JoinVertical(lipgloss.Left, emailInbox, emailContent)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, emailInbox, emailContent)
}

func EmailContent(email Email, width int) string {
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
