package messages

import (
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"time"
)

type EventDescriptionFetched struct {
	Details         models.EventDetails
	ProviderOptions []options.ProviderOption
}

type EventAsciiFetched struct {
	Ascii string
}

type EventsFetched struct {
	Events []models.Event
	Time   time.Time
}

type FetchesDone struct{}
