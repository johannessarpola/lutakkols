package fetch

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch/selectors"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/pipes"
)

type asyncSource struct{}

// Async namespaced methods for fetch
var Async asyncSource

// Events loads the events and has a rate limiting functionality for the output channel
func (a asyncSource) Events(url string, max int, ctx context.Context) chan models.Event {
	out := make(chan models.Event)

	go func() {
		defer close(out)
		ord := 0
		c := newCollector()
		var events []models.Event

		c.OnHTML(selectors.Events, func(e *colly.HTMLElement) {

			if len(events) == max && max != 0 {
				return
			}

			evt, _ := handleEvent(ord, e)
			ord++
			// Buffer to a array so that this doesn't have to wait for consumer for chan
			events = append(events, evt)

		})

		e := c.Visit(url)
		if e != nil {
			panic(e)
		}
		logger.Log.Debugf("Forwarding %d events into channel", len(events))

		// Forward the events to channel
		for _, evt := range events {
			select {
			case <-ctx.Done():
				logger.Log.Warnf("Context cancelled")
				return
			case out <- evt:
				// Sent successfully
			}
		}
	}()

	return out
}

// Images gets a channel of ascii images for event details, respecting context
// pointers are used so that there's no copying by value
func (a asyncSource) Images(eds <-chan models.EventDetails, ctx context.Context) <-chan pipes.Result[models.EventAscii] {
	out := make(chan pipes.Result[models.EventAscii])

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ed, ok := <-eds:
				if !ok {
					return
				}
				v, err := Sync.EventImage(ed.ImageURL(), ed.ID())

				var result pipes.Result[models.EventAscii]
				if err != nil {
					pipes.SendOrDone(ctx, result.WithError(err), out)
				} else {
					pipes.SendOrDone(ctx, result.WithValue(v), out)
				}
			}
		}
	}()

	return out
}

// Details gets a channel of event details for a event stream
func (a asyncSource) Details(ctx context.Context, eps <-chan models.Event) <-chan pipes.Result[models.EventDetails] {
	out := make(chan pipes.Result[models.EventDetails])
	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case ep, ok := <-eps:
				if !ok {
					return
				}
				v, err := Sync.EventDetails(ep.EventURL(), ep.ID())

				var result pipes.Result[models.EventDetails]
				if err != nil {
					pipes.SendOrDone(ctx, result.WithError(err), out)
				} else {
					pipes.SendOrDone(ctx, result.WithValue(v), out)
				}
			}
		}
	}()
	return out
}
