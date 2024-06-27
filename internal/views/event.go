package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johannessarpola/go-lutakko-gigs/internal/browser"
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/cmd"
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/constants"
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/help"
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/messages"
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/spinner"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/models"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/options"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/provider"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/logger"
	"strings"
	"time"
)

type EventViev struct {
	spinner     spinner.Model
	details     *models.EventDetails
	ascii       string
	title       string
	viewport    viewport.Model
	help        help.Model
	keyMap      EventViewKeymap
	viewReady   bool
	useWide     bool
	eventLink   string
	eventID     string
	provider    provider.Provider
	loading     bool
	loadStarted time.Time
	DataUpdated time.Time
}

func (m EventViev) Init() tea.Cmd {
	return nil
}

func viewTitle(event models.Event) string {
	return fmt.Sprintf("%s | %s", event.Headline, event.Date)
}

func InitEventView(event models.Event, provider provider.Provider) EventViev {

	ev := EventViev{
		loadStarted: time.Now(),
		DataUpdated: time.Now(),
		spinner:     newSpinner(),
		details:     &models.EventDetails{},
		title:       viewTitle(event),
		viewport:    viewport.Model{},
		useWide:     false,
		eventLink:   event.EventURL(),
		provider:    provider,
		eventID:     event.ID(),
		loading:     true,
		help:        help.New(),
		keyMap:      EventViewKeymap{},
	}

	configureView(constants.WindowSize, &ev)
	return ev
}
func (m EventViev) Ascii() string {
	return m.ascii
}

func (m EventViev) ImageLink() string {
	return m.details.ImageLink
}

func (m EventViev) Description() []string {
	return m.details.Description
}

func wideDescription(description []string) string {
	var blocks []string
	for _, block := range description {
		pp := wideDescriptionParagraphStyle.Render(block)
		blocks = append(blocks, pp)
	}
	return lipgloss.JoinVertical(lipgloss.Top, blocks...)
}

func wideASCII(ascii string) string {
	return wideAsciiStyle.Render(ascii)
}

func narrowASCII(ascii string) string {
	return fullPageAscii().Render(ascii)
}

func narrowDescription(description []string) string {
	var blocks []string
	style := fullPage()
	for _, block := range description {
		pp := style.Render(block)
		blocks = append(blocks, pp)
	}
	return lipgloss.JoinVertical(lipgloss.Top, blocks...)
}

func updateViewportContent(m *EventViev) {
	var renderedContent string

	if !m.useWide {
		na := narrowASCII(m.Ascii())
		nd := narrowDescription(m.Description())
		renderedContent = lipgloss.JoinVertical(lipgloss.Top, na, nd)
	} else {
		wa := wideASCII(m.Ascii())
		wd := wideDescription(m.Description())
		renderedContent = lipgloss.JoinHorizontal(lipgloss.Top, wa, wd)
	}
	m.viewport.SetContent(renderedContent)
}

func configureView(msg tea.WindowSizeMsg, m *EventViev) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
	m.viewport.YPosition = headerHeight
	m.viewReady = true

	if msg.Width > magicWidth {
		m.useWide = true
	} else {
		m.useWide = false
	}

}

func (m EventViev) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		c  tea.Cmd
		cs []tea.Cmd
	)
	switch msg := msg.(type) {

	case messages.EventDescriptionFetched:
		m.details = msg.Details
		gaCmd := cmd.GetAscii(msg.Details.EventID, msg.Details.ImageLink, m.provider, msg.ProviderOptions...)
		cs = append(cs, gaCmd)
	case messages.EventAsciiFetched:
		m.ascii = msg.Ascii
		doneCmd := func() tea.Msg {
			return messages.FetchesDone{}
		}
		cs = append(cs, doneCmd)
	case messages.FetchesDone:
		m.loading = false
		m.DataUpdated = time.Now()
	case tea.WindowSizeMsg:
		configureView(msg, &m)
	case spinner.TickMsg:
		m.spinner, c = m.spinner.Update(msg)
		cs = append(cs, c)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "g", "right":
			err := browser.Open(m.eventLink)
			if err != nil {
				logger.Log.Errorf("Error opening browser: %s", err.Error())
			}
		case "r", "f5":
			if m.DataUpdated.Before(time.Now().Add(-30 * time.Second)) {
				m.loading = true
				cs = append(cs, m.Refresh())
			} else {
				logger.Log.Debug("ignoring refresh")
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "backspace", "left":
			return initializeList(m.provider)
		}
	default:
		return m, nil
	}

	m.viewport, c = m.viewport.Update(msg)
	cs = append(cs, c)
	updateViewportContent(&m) // refresh every update

	return m, tea.Batch(cs...)
}

func (m EventViev) View() string {
	if m.loading {
		w := m.viewport.Width
		h := m.viewport.Height
		placedSpin := lipgloss.Place(w, h, 0.5, 0.5, m.spinner.View())
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), placedSpin, m.footerView())
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m EventViev) headerView() string {
	title := singleTitleStyle.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m EventViev) GetUpdatedAt() string {
	if m.details != nil {
		dataUpdated := m.details.UpdatedAt.Format("2006-01-02 15:04:05")
		return fmt.Sprintf("updated at %s", dataUpdated)
	}
	return ""
}

func (m EventViev) footerView() string {
	uts := updatedAtStyle.Render(m.GetUpdatedAt())
	scrollPercent := fmt.Sprintf("%3.f%% ", m.viewport.ScrollPercent()*100)
	contentBlock := lipgloss.JoinHorizontal(lipgloss.Top, scrollPercent, uts)

	infoBox := infoBoxStyle.Render(contentBlock)
	hp := m.help.View(m.keyMap)

	w := m.viewport.Width

	r := lipgloss.PlaceHorizontal(w/2, lipgloss.Right, infoBox)
	l := lipgloss.PlaceHorizontal(w/2, lipgloss.Left, hp)
	return footerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, l, r))
}

func (m EventViev) Refresh() tea.Cmd {
	// chains ot ascii fetch as well
	updateDetails := cmd.GetDetails(m.eventID, m.eventLink, m.provider, options.SkipCache)
	return updateDetails
}

func setupEventView(event models.Event, provider provider.Provider) (tea.Model, tea.Cmd) {
	eventView := InitEventView(event, provider)
	getDetailsCmd := cmd.GetDetails(event.ID(), event.EventLink, provider)
	_, updateCmd := eventView.Update(constants.WindowSize)
	return eventView, tea.Batch(eventView.spinner.Tick, updateCmd, getDetailsCmd)
}
