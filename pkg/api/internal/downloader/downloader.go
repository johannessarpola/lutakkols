// Package downloader contains the downloader implementation which is used with `sync` when updating offline data
package downloader

import (
	"errors"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/lutakkols/pkg/workset"
	"github.com/johannessarpola/lutakkols/pkg/writer"
	"time"
)

type Downloader struct {
	syncPoolSize     int
	onlineURL        string
	eventsPath       string
	eventDetailsPath string
	defaultOpts      []options.ProviderOption
}

func validateParameters(onlineURL string, eventsPath string, eventDetailsPath string, syncPoolSize int) bool {
	if len(onlineURL) == 0 {
		logger.Log.Errorf("Misconfiguration: events source URL is empty")
		return false
	}

	if len(eventsPath) == 0 {
		logger.Log.Errorf("Misconfiguration: Events path is empty")
		return false
	}

	if len(eventDetailsPath) == 0 {
		logger.Log.Errorf("Misconfiguration: eventDetails path is empty")
		return false
	}

	if syncPoolSize <= 0 {
		logger.Log.Errorf("Misconfiguration: syncPoolSize is less than 0")
		return false
	}

	return true
}

// New instantiate the offline providers, the onlineURL will only be used if the data is synchronized on demand
func New(
	onlineURL string,
	eventsPath string,
	eventDetailsPath string,
	syncPoolSize int,
	opts ...options.ProviderOption,
) (*Downloader, error) {

	if validateParameters(onlineURL, eventsPath, eventDetailsPath, syncPoolSize) == false {
		return nil, errors.New("invalid parameters")
	}

	return &Downloader{
		syncPoolSize:     syncPoolSize,
		onlineURL:        onlineURL,
		eventsPath:       eventsPath,
		eventDetailsPath: eventDetailsPath,
		defaultOpts:      opts,
	}, nil

}

// DownloadData updates the data written on the disk
func (m *Downloader) Download() error {
	if len(m.onlineURL) == 0 {
		return errors.New("no online url provided")
	}
	events, err := downloadEvents(m.onlineURL, m.eventsPath)
	if err != nil {
		return err
	}

	if events == nil {
		return errors.New("no events found")
	}

	err = downloadDetails(events, m.eventDetailsPath, m.syncPoolSize)
	if err != nil {
		return err
	}
	return nil
}

// downloadEvents downloads the event list from an url and writes it into a json file
func downloadEvents(sourceURL string, outPath string) ([]models.Event, error) {
	events, err := fetch.Events(sourceURL)
	if err != nil {
		return nil, err
	}

	err = writer.WriteJson(events, outPath, writer.PrettyPrint)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("GetEvents written to file successfully to %s from %s.\n", outPath, sourceURL)
	return events, nil
}

// downloadDetails downloads all details for the array of events concurrently and writes them into a single file
func downloadDetails(events []models.Event, outPath string, concurrentSize int) error {
	var (
		rs []*models.EventDetails
	)

	var jobs []workset.Task[*models.EventDetails]
	for _, event := range events {
		detailsFetch := func() (*models.EventDetails, error) {
			e := event // capture variable
			details, err := fetch.EventDetails(e.EventLink, e.ID())
			if err != nil {
				logger.Log.Errorf("could not fetch details %s", err.Error())
				return nil, err
			}
			logger.Log.Debugf("Downloaded details for %s", e.Headline)
			return details, nil
		}
		jobs = append(jobs, detailsFetch)
	}

	timeout := time.Second * 10
	ws := workset.NewWorkSet(jobs, concurrentSize, timeout)
	for _, result := range ws.Collect() {
		if result.Value != nil {
			rs = append(rs, result.Value)
		} else {
			logger.Log.Errorf("could not fetch event details %s", result.Error)
		}
	}
	err := writer.WriteJson(rs, outPath, writer.PrettyPrint)
	if err != nil {
		return err
	}
	logger.Log.Debugf("%d event details written to %s successfully.\n", len(rs), outPath)

	return nil
}
