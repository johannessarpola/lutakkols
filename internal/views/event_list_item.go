package views

import (
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"strings"
)

type EventViewListItem struct {
	Event models.Event
}

func (i EventViewListItem) Title() string { return i.Event.Headline }
func (i EventViewListItem) Description() string {
	sb := strings.Builder{}
	sb.WriteString(i.Event.Date)
	sb.WriteString(" · ")

	for _, bp := range i.Event.BulletPoints {
		sb.WriteString(bp)
		sb.WriteString(" · ")
	}

	return sb.String()
}
func (i EventViewListItem) FilterValue() string { return i.Event.Headline }
