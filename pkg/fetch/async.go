package fetch

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch/selectors"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"sync/atomic"
	"time"
)

type Source interface {
	Events(url string) (<-chan Result[*models.Event], <-chan error)
	Images(eds <-chan models.EventDetailsPartial) <-chan Result[*models.EventAscii]
	Details(eps <-chan models.EventPartial) <-chan Result[*models.EventDetails]
}

type Result[T any] struct {
	Val T
	Err error
}

func (r Result[T]) Value() T {
	return r.Val
}

func (r Result[T]) Error() string {
	return r.Err.Error()
}

type AsyncSource struct{}

// Events loads the events and has a rate limiting functionality for the output channel
func (a AsyncSource) Events(url string, waitTime time.Duration, context context.Context) (<-chan Result[*models.Event], <-chan error) {
	rateLimit := time.NewTicker(waitTime)

	resChan := make(chan Result[*models.Event])
	errChan := make(chan error, 1)

	go func() {
		defer close(resChan)
		defer close(errChan)
		c := newCollector()

		var ord atomic.Int32
		ord.Store(0)
		c.OnHTML(selectors.Events, func(e *colly.HTMLElement) {
			n := time.Now()
			evt := handleEvent(&ord, e)
			r := Result[*models.Event]{
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
			errChan <- e
			return
		}
	}()
	return resChan, errChan
}

func (a AsyncSource) Images(eds <-chan models.EventDetails, context context.Context) <-chan Result[*models.EventAscii] {
	resChan := make(chan Result[*models.EventAscii])

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
				resChan <- Result[*models.EventAscii]{
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

func (a AsyncSource) Details(eps <-chan models.Event, ctx context.Context) <-chan Result[*models.EventDetails] {
	resChan := make(chan Result[*models.EventDetails])
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
				resChan <- Result[*models.EventDetails]{
					Val: v,
					Err: err,
				}
			}
		}
	}()
	return resChan
}
