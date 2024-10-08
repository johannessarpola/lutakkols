package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johannessarpola/lutakkols/internal/browser"
	"github.com/johannessarpola/lutakkols/internal/views/cmd"
	"github.com/johannessarpola/lutakkols/internal/views/constants"
	"github.com/johannessarpola/lutakkols/internal/views/help"
	"github.com/johannessarpola/lutakkols/internal/views/messages"
	"github.com/johannessarpola/lutakkols/internal/views/spinner"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/api/provider"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"time"
)

type WindowSize struct {
	w int
	h int
}

type EventList struct {
	list        list.Model
	help        help.Model
	loading     bool
	spinner     spinner.Model
	Quitting    bool
	provider    provider.Provider
	DataUpdated time.Time
	WindowSize  WindowSize
}

func massageItems(events []models.Event) []list.Item {
	items := make([]list.Item, len(events))
	for i, v := range events {
		items[i] = EventViewListItem{
			Event: v,
		}
	}
	return items
}

func setupKeybinds(m *list.Model) {

	m.KeyMap = eventListKeymap()

	m.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			GoToEventPage(),
			RefreshPage(),
		}
	}
	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			GoToEventPage(),
			RefreshPage(),
		}
	}
}

func newHelp() help.Model {
	m := help.New()
	return m
}

func setupListModel(slm *list.Model) {
	setupKeybinds(slm)
	slm.Title = constants.Title
	slm.SetShowHelp(false) // TODO customize sometime
	slm.SetShowStatusBar(false)
	slm.SetShowTitle(false)
	slm.SetFilteringEnabled(false) // TODO Fix sometime
	slm.Styles.PaginationStyle = paginationStyle
	slm.Styles.HelpStyle = footerStyle
}

func (m EventList) configureList(items []list.Item) list.Model {

	headerHeight := lipgloss.Height(m.Header())
	footerHeight := lipgloss.Height(m.Footer())
	delegate := customizedDelegate()

	slm := list.New(items, delegate, constants.WindowSize.Width, constants.WindowSize.Height-headerHeight-footerHeight)
	setupListModel(&slm)
	return slm
}

func emptyList() list.Model {
	delegate := list.NewDefaultDelegate()
	slm := list.New(make([]list.Item, 0), delegate, defaultWidth, defaultHeight)
	setupListModel(&slm)
	return slm
}

func NewEventsList(provider provider.Provider) EventList {

	return EventList{
		Quitting:    false,
		loading:     true,
		spinner:     newSpinner(),
		list:        emptyList(),
		provider:    provider,
		DataUpdated: time.Now(),
		help:        newHelp(),
	}
}

func (m EventList) Init() tea.Cmd {
	ge := cmd.GetEvents(m.provider)
	return tea.Sequence(m.spinner.Tick, ge)
}

func (m EventList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var c tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, c = m.spinner.Update(msg)
		return m, c
	case tea.WindowSizeMsg:
		// we need space for the custom header
		constants.WindowSize = msg
		newHeight := msg.Height - lipgloss.Height(m.Header()) - lipgloss.Height(m.Footer())
		newMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: newHeight,
		}
		m.list, c = m.list.Update(newMsg)
		m.WindowSize = WindowSize{
			w: msg.Width,
			h: msg.Height,
		}
		return m, c
	case messages.EventsFetched:
		i := massageItems(msg.Events)
		m.DataUpdated = msg.Time
		m.list = m.configureList(i)
		m.loading = false
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "g":
			selectedEvent := m.list.SelectedItem().(EventViewListItem)
			err := browser.Open(selectedEvent.Event.EventLink)
			if err != nil {
				panic(err)
			}
		case "r", "f5":
			if m.DataUpdated.Before(time.Now().Add(-30 * time.Second)) {
				m.loading = true
				return m, cmd.GetEvents(m.provider, options.SkipCache)
			} else {
				logger.Log.Debug("ignoring refresh")
			}
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			selectedEvent := m.list.SelectedItem().(EventViewListItem)
			return setupEventView(selectedEvent.Event, m.provider)
		}
	}

	m.list, c = m.list.Update(msg)
	return m, c
}

func initializeList(provider provider.Provider) (tea.Model, tea.Cmd) {
	newList := NewEventsList(provider)
	init := newList.Init()
	tick := newList.spinner.Tick
	_, update := newList.Update(constants.WindowSize)
	return newList, tea.Sequence(init, update, tick)

}

func (m EventList) Footer() string {
	uts := updatedAtStyle.Render(m.GetUpdatedAt())
	padded := infoBoxStyle.Render(uts)
	h := m.help.View(m.list)
	l := lipgloss.PlaceHorizontal(constants.WindowSize.Width/2, lipgloss.Left, h)
	r := lipgloss.PlaceHorizontal(constants.WindowSize.Width/2, lipgloss.Right, padded)

	return footerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, l, r))
}

func (m EventList) GetUpdatedAt() string {
	dataUpdated := m.DataUpdated.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("updated at %s", dataUpdated)
}

func (m EventList) Header() string {
	title := constants.Title
	r1 := titleTextStyle.Render(title)
	return titleBoxStyle.Render(r1)
}

func (m EventList) View() string {
	if m.Quitting {
		//	return quitTextStyle.Render("Quitting ...")
	}
	header := m.Header()
	footer := m.Footer()
	availableHeight := m.list.Height() - lipgloss.Height(header) - lipgloss.Height(footer)
	if m.loading {
		p := lipgloss.Place(defaultWidth, availableHeight, 0.5, 0.5, m.spinner.View())
		return lipgloss.JoinVertical(lipgloss.Top, header, p, footer)
	}

	l := m.list.View()
	return lipgloss.JoinVertical(lipgloss.Top, header, l, footer)
}
