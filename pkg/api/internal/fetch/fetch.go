// Package fetch contains the methods to extract the relevant models from the HTML from the source URL
// also handles the conversion of images to ascii art with the image2ascii library
package fetch

import (
	"bytes"
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/internal/fetch/selectors"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/qeesung/image2ascii/convert"
	"image"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// EventImage fetches normal image file and turns it into an ascii art
func EventImage(url string) (string, error) {
	img, err := downloadImage(url)
	if err != nil {
		return "", FailedFetch{err: err, url: url}
	}

	converter := convert.NewImageConverter()
	return converter.Image2ASCIIString(*img, defaultConvertorOptions()), nil

}

// Events fetches the events from the source
func Events(url string) ([]models.Event, error) {
	c := newCollector()
	var events []models.Event
	var ord atomic.Int32
	ord.Store(0)

	c.OnHTML(selectors.Events, func(e *colly.HTMLElement) {
		evt := extractEvent(e)
		evt.UpdatedAt = time.Now()
		evt.Order = ord.Add(1)
		events = append(events, evt)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, FailedFetch{err: err, url: url}
	}

	return events, nil
}

// EventDetails fetches the eventDetails for eventUrl from source
func EventDetails(url string, eventId string) (*models.EventDetails, error) {
	c := newCollector()
	ed := models.EventDetails{}
	ed.EventID = eventId
	ed.UpdatedAt = time.Now()

	// extract the product info for event
	c.OnHTML(selectors.EventProductInfo, func(e *colly.HTMLElement) {
		ed.ProductInfo = extractProdductInfo(e)
		ed.PlayTimes = extractPlayTimes(e)
	})

	// extract product summary
	c.OnHTML(selectors.EventSummary, func(e *colly.HTMLElement) {
		ed.Description = extractSummary(e)
		ed.ImageLink = extractImageLink(e)
	})

	// extract tickets
	c.OnHTML(selectors.EventTickets, func(e *colly.HTMLElement) {
		ed.Tickets = extractTicketPrices(e)
	})

	// extract door price
	c.OnHTML(selectors.DoorPrice, func(e *colly.HTMLElement) {
		ed.DoorPrice = extractDoorPrice(e)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, FailedFetch{err: err, url: url}
	}
	return &ed, nil

}

func downloadImage(url string) (*image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func defaultConvertorOptions() *convert.Options {
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 40
	convertOptions.FixedHeight = 20
	convertOptions.StretchedScreen = false
	return &convertOptions
}
