// Package online contains the main provider to use with internet connection
package online

import (
	"errors"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/caching"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"time"
)

type Provider struct {
	sourceURL   string
	fetchCache  *caching.EventCache
	defaultOpts []options.ProviderOption
}

var ttlOptions = caching.TTLOptions{
	EventsTTL:       time.Duration(5) * time.Minute,
	EventDetailsTTL: time.Duration(5) * time.Minute,
	EventAsciiTTL:   time.Duration(30) * time.Minute,
	DefaultTTL:      time.Duration(5) * time.Minute,
	Capacity:        1000,
}

func (m *Provider) useCache(opts []options.ProviderOption) bool {
	return !options.Has(options.SkipCache, m.withInitialOpts(opts)) && m.fetchCache != nil
}

func New(eventsSourceURL string, opts ...options.ProviderOption) Provider {

	c, err := caching.New(ttlOptions)
	if err != nil {
		// we can operate without cache
		logger.Log.Warnf("Err initializing cache: %v", err)
	}

	return Provider{
		sourceURL:   eventsSourceURL,
		fetchCache:  c,
		defaultOpts: opts,
	}
}

func (m *Provider) withInitialOpts(additionalOpts []options.ProviderOption) []options.ProviderOption {
	return append(m.defaultOpts, additionalOpts...)
}

func (m *Provider) GetAscii(eventID string, imageURL string, opts ...options.ProviderOption) (models.EventAscii, error) {
	var ea models.EventAscii
	var err error

	if len(imageURL) == 0 {
		return ea, errors.New("image link missing")
	}
	if m.useCache(opts) {
		value, ts, ok := m.fetchCache.GetAscii(eventID)
		if ok {
			value.UpdatedAt = ts
			logger.Log.Debugf("fetched ea from caching with id %s", eventID)
			return value, nil
		}
	}

	ea, err = fetch.Sync.EventImage(imageURL, eventID)
	if err == nil {
		m.fetchCache.SetAscii(eventID, ea)
	}
	return ea, err
}

func (m *Provider) GetDetails(eventID string, eventURL string, opts ...options.ProviderOption) (models.EventDetails, error) {
	var ed models.EventDetails
	var err error

	if len(eventURL) == 0 {
		return ed, errors.New("event url missing")
	}
	if m.useCache(opts) {
		value, ts, ok := m.fetchCache.GetDetails(eventID)
		if ok {
			value.UpdatedAt = ts
			logger.Log.Debugf("fetched details from caching with id %s", eventID)
			return value, nil
		}
	}

	ed, err = fetch.Sync.EventDetails(eventURL, eventID)
	if err == nil {
		m.fetchCache.SetDetails(eventID, ed)
	}

	return ed, err
}

func (m *Provider) GetEvents(opts ...options.ProviderOption) (*models.Events, error) {
	if m.useCache(opts) {
		value, ts, ok := m.fetchCache.GetEvents()
		if ok {
			logger.Log.Debugf("fetched from caching events with timestamp %s\n", ts.String())
			return value, nil
		}
	}

	list, err := fetch.Sync.Events(m.sourceURL)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("no events found")
	}

	events := models.Events{
		Events:    list,
		UpdatedAt: time.Now(),
	}

	if list != nil {
		m.fetchCache.SetEvents(events)
	}

	return &events, nil

}
