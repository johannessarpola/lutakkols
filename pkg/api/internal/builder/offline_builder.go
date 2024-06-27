package builder

import (
	"errors"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/offline"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/logger"
)

type OfflineBuilder struct {
	DefaultOpts        []options.ProviderOption
	EventSourceFsPath  string
	EventDetailsFsPath string
	AsciiGen           func(string, string) string
}

func (b *OfflineBuilder) WithDefaultOpts(opts ...options.ProviderOption) *OfflineBuilder {
	b.DefaultOpts = opts
	return b
}

func (b *OfflineBuilder) WithEventSourceFsPath(path string) *OfflineBuilder {
	b.EventSourceFsPath = path
	return b
}

func (b *OfflineBuilder) WithEventDetailsFsPath(path string) *OfflineBuilder {
	b.EventDetailsFsPath = path
	return b
}

func (b *OfflineBuilder) WithAsciiGen(genFunc func(string, string) string) *OfflineBuilder {
	b.AsciiGen = genFunc
	return b
}

func (b *OfflineBuilder) validateParameters() bool {
	if len(b.EventSourceFsPath) == 0 {
		logger.Log.Error("Misconfiguration: events path is empty")
		return false
	}

	if len(b.EventDetailsFsPath) == 0 {
		logger.Log.Error("Misconfiguration: eventDetails path is empty")
		return false
	}

	if b.AsciiGen == nil {
		logger.Log.Error("Misconfiguration: AsciiGen function is nil")
		return false
	}

	return true
}

func (b *OfflineBuilder) Build() (*offline.Provider, error) {

	if !b.validateParameters() {
		return nil, errors.New("invalid parameters")
	}

	p := offline.New(b.EventSourceFsPath,
		b.EventDetailsFsPath,
		b.AsciiGen,
		b.DefaultOpts...)
	return &p, nil

}
