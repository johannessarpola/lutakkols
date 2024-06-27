package spinner

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"sync"
	"time"
)

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var (
	lastID int
	idMtx  sync.Mutex
)

// Return the next ID we should use on the Model.
func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

// Spinner is a set of frames used in animating the spinner.
type Spinner struct {
	Frames      []string
	FPS         time.Duration
	Description string
}

var (
	LutakkoSpinner = Spinner{
		Frames:      []string{"🌍", "🌎", "🌏"},
		FPS:         time.Second / 4,
		Description: "loading",
	}
)

// Model contains the state for the spinner. Use New to create new models
// rather than using Model as a struct literal.
type Model struct {
	Spinner Spinner
	Style   lipgloss.Style

	frame int
	id    int
	tag   int
}

// ID returns the spinner's unique ID.
func (m Model) ID() int {
	return m.id
}

// New returns a model with default values.
func New(opts ...Option) Model {
	m := Model{
		Spinner: LutakkoSpinner,
		id:      nextID(),
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// TickMsg indicates that the timer has ticked and we should render a frame.
type TickMsg struct {
	Time time.Time
	tag  int
	ID   int
}

// Update is the Tea update function.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		// If an ID is set, and the ID doesn't belong to this spinner, reject
		// the message.
		if msg.ID > 0 && msg.ID != m.id {
			return m, nil
		}

		// If a tag is set, and it's not the one we expect, reject the message.
		// This prevents the spinner from receiving too many messages and
		// thus spinning too fast.
		if msg.tag > 0 && msg.tag != m.tag {
			return m, nil
		}

		m.frame++
		if m.frame >= len(m.Spinner.Frames) {
			m.frame = 0
		}

		m.tag++
		return m, m.tick(m.id, m.tag)
	default:
		return m, nil
	}
}

// View renders the model's view.
func (m Model) View() string {

	if m.frame >= len(m.Spinner.Frames) {
		return "(error)"
	}

	str := fmt.Sprintf("%s %s%s", m.Spinner.Frames[m.frame], m.Spinner.Description, strings.Repeat(".", m.frame))

	return m.Style.Render(str)
}

// Tick is the command used to advance the spinner one frame. Use this command
// to effectively start the spinner.
func (m Model) Tick() tea.Msg {
	return TickMsg{
		// The time at which the tick occurred.
		Time: time.Now(),

		// The ID of the spinner that this message belongs to. This can be
		// helpful when routing messages, however bear in mind that spinners
		// will ignore messages that don't contain ID by default.
		ID: m.id,

		tag: m.tag,
	}
}

func (m Model) tick(id, tag int) tea.Cmd {
	return tea.Tick(m.Spinner.FPS, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

// Option is used to set options in New. For example:
//
//	spinner := New(WithSpinner(Dot))
type Option func(*Model)

// WithSpinner is an option to set the spinner.
func WithSpinner(spinner Spinner) Option {
	return func(m *Model) {
		m.Spinner = spinner
	}
}

// WithStyle is an option to set the spinner style.
func WithStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.Style = style
	}
}
