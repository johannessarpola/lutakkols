package views

import "github.com/charmbracelet/lipgloss"

func GenerateOfflineAscii(_ string, url string) string {
	na := "asici not available in offline mode:"

	block := lipgloss.JoinVertical(lipgloss.Top, na, url)
	asc := lipgloss.Place(asciiWidth, asciiHeight, lipgloss.Center, lipgloss.Center, asciiPlaceholderStyle.Render(block), lipgloss.WithWhitespaceChars("."))
	return asc
}
