package builder

import (
	"errors"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/internal/downloader"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/logger"
)

const defaultPoolSize = 4

type DownloaderBuilder struct {
	EventSourceFsPath  string
	EventDetailsFsPath string
	EventsSourceURL    string
	SyncPoolSize       int
}

func (b *DownloaderBuilder) WithEventSourceFsPath(path string) *DownloaderBuilder {
	b.EventSourceFsPath = path
	return b
}

func (b *DownloaderBuilder) WithEventDetailsFsPath(path string) *DownloaderBuilder {
	b.EventDetailsFsPath = path
	return b
}

func (b *DownloaderBuilder) WitEventsSourceURL(path string) *DownloaderBuilder {
	b.EventsSourceURL = path
	return b
}

func (b *DownloaderBuilder) WitSyncPoolSize(size int) *DownloaderBuilder {
	b.SyncPoolSize = size
	return b
}

func (b *DownloaderBuilder) validateParameters() bool {
	if len(b.EventSourceFsPath) == 0 {
		logger.Log.Error("Misconfiguration: events path is empty")
		return false
	}

	if len(b.EventDetailsFsPath) == 0 {
		logger.Log.Error("Misconfiguration: eventDetails path is empty")
		return false
	}

	if len(b.EventsSourceURL) == 0 {
		logger.Log.Error("Misconfiguration: events source URL is empty")
		return false
	}

	return true
}

// Build constructs obj which can be used to download data
func (b *DownloaderBuilder) Build() (*downloader.Downloader, error) {
	if !b.validateParameters() {
		return nil, errors.New("invalid parameters")
	}

	if b.SyncPoolSize == 0 {
		b.SyncPoolSize = defaultPoolSize
	}

	return downloader.New(b.EventsSourceURL, b.EventSourceFsPath, b.EventDetailsFsPath, b.SyncPoolSize)
}
