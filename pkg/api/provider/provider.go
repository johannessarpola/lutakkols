package provider

import (
	"errors"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/builder"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
)

type Config struct {
	EventsSourceURL    string
	DefaultOpts        []options.ProviderOption
	EventSourceFsPath  string
	EventDetailsFsPath string
	AsciiGen           func(string, string) string
}

type Provider interface {
	GetEvents(opts ...options.ProviderOption) (*models.Events, error)
	GetAscii(eventID string, imageURL string, opts ...options.ProviderOption) (models.EventAscii, error)
	GetDetails(eventID string, eventURL string, opts ...options.ProviderOption) (models.EventDetails, error)
}

type Downloader interface {
	Download() error
}

// New constructs the correct provider from the configuration
func New(config *Config, opt options.TypeOption) (Provider, error) {
	switch opt {
	case options.UseOnline:
		b := (&builder.OnlineBuilder{}).
			WitEventsSourceURL(config.EventsSourceURL).
			WithDefaultOpts(config.DefaultOpts...)
		return b.Build()
	case options.UseOffline:
		b := (&builder.OfflineBuilder{}).
			WithEventSourceFsPath(config.EventSourceFsPath).
			WithDefaultOpts(config.DefaultOpts...).
			WithEventDetailsFsPath(config.EventDetailsFsPath).
			WithAsciiGen(config.AsciiGen)
		return b.Build()
	default:
		return nil, errors.New("invalid option")
	}
}
