package builder

import (
	"errors"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/online"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/logger"
)

type OnlineBuilder struct {
	EventsSourceURL string
	DefaultOpts     []options.ProviderOption
}

func (b *OnlineBuilder) WithDefaultOpts(opts ...options.ProviderOption) *OnlineBuilder {
	b.DefaultOpts = opts
	return b
}

func (b *OnlineBuilder) WitEventsSourceURL(path string) *OnlineBuilder {
	b.EventsSourceURL = path
	return b
}

func (b *OnlineBuilder) validateParameters() bool {
	if len(b.EventsSourceURL) == 0 {
		logger.Log.Error("Misconfiguration: Event source URL is empty")
		return false
	}

	return true
}

func (b *OnlineBuilder) Build() (*online.Provider, error) {
	if !b.validateParameters() {
		return nil, errors.New("invalid parameters")
	}

	p := online.New(b.EventsSourceURL, b.DefaultOpts...)
	return &p, nil
}
