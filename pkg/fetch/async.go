package fetch

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch/selectors"
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

type AsyncSource struct{}

func (a AsyncSource) Events(url string, context context.Context) (<-chan Result[*models.Event], <-chan error) {
	resChan := make(chan Result[*models.Event])
	errChn := make(chan error, 1)

	go func() {
		defer close(resChan)
		defer close(errChn)
		c := newCollector()

		var ord atomic.Int32
		ord.Store(0)
		c.OnHTML(selectors.Events, func(e *colly.HTMLElement) {
			evt := handleEvent(&ord, e)
			r := Result[*models.Event]{
				Val: &evt,
				Err: nil,
			}

			select {
			case <-context.Done():
				return
			default:
				fmt.Printf("sending event %s\n", r.Val.Headline)
				resChan <- r
			}

		})

		fmt.Println("visiting ", url)
		e := c.Visit(url)
		if e != nil {
			fmt.Println("error ", e)
			errChn <- e
			return
		}
	}()
	return resChan, errChn
}

func (a AsyncSource) Images(eds <-chan models.EventDetailsPartial, context context.Context) <-chan Result[*models.EventAscii] {
	resChan := make(chan Result[*models.EventAscii])

	go func() {
		defer close(resChan)
		for {
			select {
			case <-context.Done():
				return
			case ed, ok := <-eds:
				fmt.Printf("fetching image from %s\n", ed.ImageURL())
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

func (a AsyncSource) Details(eps <-chan models.EventPartial, ctx context.Context) <-chan Result[*models.EventDetails] {
	resChn := make(chan Result[*models.EventDetails])
	go func() {
		defer close(resChn)

		for {
			select {
			case <-ctx.Done():
				return
			case ep, ok := <-eps:
				if !ok {
					return
				}
				fmt.Printf("fetching details from %s\n", ep.EventURL())
				v, err := EventDetails(ep.EventURL(), ep.ID())
				fmt.Printf("sending details %s\n", ep.EventURL())
				resChn <- Result[*models.EventDetails]{
					Val: v,
					Err: err,
				}
				fmt.Println("sent event details")
			}
		}
	}()

	return resChn
}
