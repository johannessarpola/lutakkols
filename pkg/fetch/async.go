package fetch

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch/selectors"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/lutakkols/pkg/pipes"
	"sync/atomic"
	"time"
)

type AsyncSource struct{}

// Events loads the events and has a rate limiting functionality for the output channel, respecting context
// pointers are used so that there's no copying by value
func (a AsyncSource) Events(url string, waitTime time.Duration, context context.Context) <-chan pipes.Result[*models.Event] {
	rateLimit := time.NewTicker(waitTime)
	resChan := make(chan pipes.Result[*models.Event])

	go func() {
		defer close(resChan)
		c := newCollector()

		var ord atomic.Int32
		ord.Store(0)
		c.OnHTML(selectors.Events, func(e *colly.HTMLElement) {
			n := time.Now()
			evt := handleEvent(&ord, e)
			r := pipes.Result[*models.Event]{
				Val: &evt,
				Err: nil,
			}

			select {
			case <-context.Done():
				return
			case <-rateLimit.C:
				logger.Log.Debugf("Tick delay %s", time.Now().Sub(n))
				resChan <- r
			}

		})

		e := c.Visit(url)
		if e != nil {
			r := pipes.Result[*models.Event]{
				Val: nil,
				Err: e,
			}
			resChan <- r
			return
		}
	}()
	return resChan
}

// Images gets a channel of ascii images for event details, respecting context
// pointers are used so that there's no copying by value
func (a AsyncSource) Images(eds <-chan models.EventDetails, context context.Context) <-chan pipes.Result[*models.EventAscii] {
	resChan := make(chan pipes.Result[*models.EventAscii])

	go func() {
		defer close(resChan)
		for {
			select {
			case <-context.Done():
				return
			case ed, ok := <-eds:
				if !ok {
					return
				}
				v, err := EventImage(ed.ImageURL())
				resChan <- pipes.Result[*models.EventAscii]{
					Val: &models.EventAscii{
						Ascii:     v,
						EventID:   ed.ID(),
						UpdatedAt: time.Now(),
					},
					Err: err,
				}
			}
		}
	}()

	return resChan
}

// Details gets a channel of event details for a event stream, respecting context
// pointers are used so that there's no copying by value
func (a AsyncSource) Details(eps <-chan *models.Event, ctx context.Context) <-chan pipes.Result[*models.EventDetails] {
	resChan := make(chan pipes.Result[*models.EventDetails])
	go func() {
		defer close(resChan)

		for {
			select {
			case <-ctx.Done():
				return
			case ep, ok := <-eps:
				if !ok {
					return
				}
				v, err := EventDetails(ep.EventURL(), ep.ID())
				resChan <- pipes.Result[*models.EventDetails]{
					Val: v,
					Err: err,
				}
			}
		}
	}()
	return resChan
}
