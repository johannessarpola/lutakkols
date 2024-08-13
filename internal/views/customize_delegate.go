package views

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func customizeDelegate(d list.DefaultDelegate) list.ItemDelegate {
	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#4242f5", Dark: "#f5d742"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#4242f5", Dark: "#f5d742"}).
		Padding(0, 0, 0, 1)

	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#4242f5", Dark: "#f5d742"})

	return d
}
