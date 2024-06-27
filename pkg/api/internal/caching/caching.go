package caching

import (
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/internal/cache"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/models"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/logger"
	"time"
)

type TTLOptions struct {
	EventsTTL       time.Duration
	EventDetailsTTL time.Duration
	EventAsciiTTL   time.Duration
	DefaultTTL      time.Duration
	Capacity        int
}

func (t TTLOptions) ttlOpts() []cache.ItemTTL {
	events := cache.ItemTTL{
		Prefix: prefixKey(models.Events{}),
		TTL:    t.EventsTTL,
	}

	details := cache.ItemTTL{
		Prefix: prefixKey(models.EventDetails{}),
		TTL:    t.EventDetailsTTL,
	}

	ascii := cache.ItemTTL{
		Prefix: prefixKey(models.EventAscii{}),
		TTL:    t.EventAsciiTTL,
	}

	return []cache.ItemTTL{
		events, details, ascii,
	}
}

func prefixKey(a any) string {
	switch a.(type) {
	case models.Events, *models.Event:
		return "events"
	case models.EventAscii, *models.EventAscii:
		return "ascii"
	case models.EventDetails, *models.EventDetails:
		return "details"
	default:
		logger.Log.Warn("invalid cacheKey in keyFunc")
		return "unknown"
	}
}

const eventsKey = "all"

type EventCache struct {
	internalCache *cache.Cache
}

func (f EventCache) Clear() {
	f.internalCache.Clear()
}

func (f EventCache) GetEvents() (*models.Events, time.Time, bool) {
	ret := models.Events{}
	v, ts, ok := f.internalCache.Get(prefixKey(ret), eventsKey)
	if ok {
		ret, ok = v.(models.Events)
		if ok {
			return &ret, ts, ok
		}
	}
	return nil, time.Time{}, false

}

func (f EventCache) GetDetails(eventID string) (*models.EventDetails, time.Time, bool) {
	ret := models.EventDetails{}
	v, ts, ok := f.internalCache.Get(prefixKey(ret), eventID)
	if ok {
		ret, ok = v.(models.EventDetails)
		if ok {
			return &ret, ts, ok
		}
	}
	return nil, time.Time{}, false
}

func (f EventCache) GetAscii(eventID string) (*models.EventAscii, time.Time, bool) {
	ret := models.EventAscii{}
	v, ts, ok := f.internalCache.Get(prefixKey(ret), eventID)
	if ok {
		ret, ok = v.(models.EventAscii)
		if ok {
			return &ret, ts, ok
		}
	}
	return nil, time.Time{}, false
}

func (f EventCache) SetDetails(eventID string, ed models.EventDetails) bool {
	return f.internalCache.Set(prefixKey(ed), eventID, ed)
}

func (f EventCache) SetAscii(eventID string, ea models.EventAscii) bool {
	return f.internalCache.Set(prefixKey(ea), eventID, ea)
}

func (f EventCache) SetEvents(events models.Events) bool {
	return f.internalCache.Set(prefixKey(events), eventsKey, events)
}

func New(options TTLOptions) (*EventCache, error) {
	ic, err := cache.New(options.Capacity, options.DefaultTTL, options.ttlOpts()...)
	if err != nil {
		return nil, err
	}
	return &EventCache{
		internalCache: ic,
	}, nil
}
