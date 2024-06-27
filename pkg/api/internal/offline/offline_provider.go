// Package offline contains the provider to use when running offline mode
package offline

import (
	"github.com/johannessarpola/lutakkols/pkg/api/internal/caching"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/loadfs"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"time"
)

type Provider struct {
	eventsPath       string
	eventDetailsPath string
	fetchCache       *caching.EventCache
	defaultOpts      []options.ProviderOption
	asciiGen         func(string, string) string
}

const singleTTL = time.Duration(120) * time.Minute

var ttlOptions = caching.TTLOptions{
	EventsTTL:       singleTTL,
	EventDetailsTTL: singleTTL,
	EventAsciiTTL:   singleTTL,
	DefaultTTL:      singleTTL,
	Capacity:        1000,
}

// New instantiate the offline providers, the onlineURL will only be used if the data is synchronized on demand
func New(
	eventsPath string,
	eventDetailsPath string,
	asciiGenerator func(string, string) string,
	opts ...options.ProviderOption,
) Provider {

	c, err := caching.New(ttlOptions)
	if err != nil {
		// we can operate without cache
		logger.Log.Warnf("could not create cache: %v", err)
	}

	return Provider{
		eventsPath:       eventsPath,
		eventDetailsPath: eventDetailsPath,
		defaultOpts:      opts,
		asciiGen:         asciiGenerator,
		fetchCache:       c,
	}
}

func (m *Provider) useCache(opts []options.ProviderOption) bool {
	return !options.Has(options.SkipCache, m.withInitialOpts(opts)) && m.fetchCache != nil
}

func (m *Provider) withInitialOpts(additionalOpts []options.ProviderOption) []options.ProviderOption {
	return append(m.defaultOpts, additionalOpts...)
}

// GetAscii just prints placeholder string since it is not possible to do this offline (yet)
func (m *Provider) GetAscii(eventID string, imageURL string, _ ...options.ProviderOption) (*models.EventAscii, error) {
	return &models.EventAscii{
		Ascii:   m.asciiGen(eventID, imageURL),
		EventID: eventID,
	}, nil
}

func (m *Provider) GetDetails(eventID string, _ string, opts ...options.ProviderOption) (*models.EventDetails, error) {
	if m.useCache(opts) {
		value, ts, ok := m.fetchCache.GetDetails(eventID)
		if ok {
			value.UpdatedAt = ts
			logger.Log.Infof("fetched ascii from caching with id %s", eventID)
			return value, nil
		}
	}

	ed, err := loadfs.EventDetails(eventID, m.eventDetailsPath)
	if err != nil {
		return nil, err
	}
	if ed != nil {
		m.fetchCache.SetDetails(eventID, *ed)
	}
	return ed, nil
}

func (m *Provider) GetEvents(opts ...options.ProviderOption) (*models.Events, error) {
	if m.useCache(opts) {
		value, ts, ok := m.fetchCache.GetEvents()
		if ok {
			logger.Log.Infof("fetched from caching events with timestamp %s\n", ts.String())
			return value, nil
		}
	}

	events, err := loadfs.Events(m.eventsPath)
	if events != nil {
		m.fetchCache.SetEvents(*events)
	}
	return events, err
}
