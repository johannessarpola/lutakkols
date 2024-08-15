package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/johannessarpola/lutakkols/internal/views/constants"
)

const defaultWidth = 100 // parameter?
const defaultHeight = 30 // parameter=
const asciiWidth = 40
const asciiHeight = 32

var (
	magicWidth = 110

	singleTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	wideDescriptionParagraphStyle = lipgloss.NewStyle().Width(70).MarginTop(1).MarginLeft(2)
	wideAsciiStyle                = lipgloss.NewStyle().Width(40)

	titleTextStyle        = lipgloss.NewStyle()
	infoBoxStyle          = lipgloss.NewStyle().Align(lipgloss.Right).PaddingRight(2)
	footerStyle           = lipgloss.NewStyle().PaddingTop(2).PaddingLeft(2).PaddingBottom(1)
	paginationStyle       = lipgloss.NewStyle().PaddingLeft(2)
	subduedColor          = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	titleBoxStyle         = lipgloss.NewStyle().MarginBottom(1)
	updatedAtStyle        = lipgloss.NewStyle().Foreground(subduedColor)
	asciiPlaceholderStyle = lipgloss.NewStyle().Width(asciiWidth * 0.7).Padding(1)
)

const magicReduce = 6

func fullPage() lipgloss.Style {
	return lipgloss.NewStyle().Width(constants.WindowSize.Width - magicReduce).Margin(1)
}

func fullPageAscii() lipgloss.Style {
	return lipgloss.NewStyle().Width(constants.WindowSize.Width - magicReduce).Align(lipgloss.Center).Margin(1)
}
