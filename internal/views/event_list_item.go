package views

import (
	"github.com/johannessarpola/lutakkols/pkg/api/models"
)

type EventViewListItem struct {
	Event models.Event
}

func (i EventViewListItem) Title() string       { return i.Event.Headline }
func (i EventViewListItem) Description() string { return i.Event.Date }
func (i EventViewListItem) FilterValue() string { return i.Event.Headline }
