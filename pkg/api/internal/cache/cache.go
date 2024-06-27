package cache

import (
	"fmt"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/logger"
	"github.com/maypok86/otter"
	"time"
)

type cachedData struct {
	CachedAt time.Time
	Data     any
}

// Cache wrapper for otter.Cache so that we can access the timestamps and have prefix->ttl based expiration
type Cache struct {
	defaultTTL time.Duration
	tm         map[string]time.Duration
	otterCache *otter.CacheWithVariableTTL[string, cachedData]
}

func (f Cache) Get(prefix string, key string) (any, time.Time, bool) {
	ck := cacheKey(prefix, key)
	v, ok := f.otterCache.Get(ck)
	return v.Data, v.CachedAt, ok
}

func cacheKey(prefix string, key string) string {
	return fmt.Sprintf("%s_%s", prefix, key)
}

func (f Cache) Set(prefix string, key string, item any) bool {
	ck := cacheKey(prefix, key)
	ttl, ok := f.tm[prefix]
	if !ok {
		ttl = f.defaultTTL
	}
	c := cachedData{
		CachedAt: time.Now(),
		Data:     item,
	}

	return f.otterCache.Set(ck, c, ttl)
}

func (f Cache) Clear() {
	f.otterCache.Clear()
}

type ItemTTL struct {
	Prefix string
	TTL    time.Duration
}

func buildCache[K comparable, V any](capacity int) (otter.CacheWithVariableTTL[K, V], error) {
	return otter.MustBuilder[K, V](capacity).
		Cost(func(key K, value V) uint32 {
			return 1
		}).
		WithVariableTTL().
		Build()
}

func setupTTLs(ttls []ItemTTL) map[string]time.Duration {
	tm := make(map[string]time.Duration)
	for _, ttl := range ttls {
		logger.Log.Debugf("Setting up TTL for %s with TTL %s", ttl.Prefix, ttl.TTL.String())
		tm[ttl.Prefix] = ttl.TTL
	}
	return tm
}

func New(capacity int, defaultTTL time.Duration, ttls ...ItemTTL) (*Cache, error) {
	tm := setupTTLs(ttls)
	cache, err := buildCache[string, cachedData](capacity)
	if err != nil {
		return nil, err
	}
	return &Cache{
		defaultTTL: defaultTTL,
		tm:         tm,
		otterCache: &cache,
	}, nil
}
