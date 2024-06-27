package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func Back() key.Binding {
	return key.NewBinding(
		key.WithKeys("backspace", "left"),
		key.WithHelp("<-/backspace", "back"),
	)
}

func GoToEventPage() key.Binding {
	return key.NewBinding(
		key.WithKeys("g", "right"),
		key.WithHelp("g/->", "browser"),
	)
}

func RefreshPage() key.Binding {
	return key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	)
}

func CursorUp() key.Binding {
	return key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	)
}

func CursorDown() key.Binding {
	return key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	)
}

// eventListKeymap returns a default set of keybindings.
func eventListKeymap() list.KeyMap {
	return list.KeyMap{
		// Browsing.
		CursorUp:   CursorUp(),
		CursorDown: CursorDown(),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "h", "pgup", "b", "u"),
			key.WithHelp("←/h/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "l", "pgdown", "f", "d"),
			key.WithHelp("→/l/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "h"),
			key.WithHelp("h/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply filter"),
		),

		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

type EventViewKeymap struct{}

func (m EventViewKeymap) ShortHelp() []key.Binding {
	group := []key.Binding{
		CursorUp(),
		CursorDown(),
		RefreshPage(),
		GoToEventPage(),
		Back(),
	}
	return group
}

func (m EventViewKeymap) FullHelp() [][]key.Binding {
	group := []key.Binding{
		CursorUp(),
		CursorDown(),
		RefreshPage(),
		GoToEventPage(),
		Back(),
	}
	return [][]key.Binding{
		group,
	}
}
