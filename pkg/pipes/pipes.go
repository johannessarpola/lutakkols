package pipes

import (
	"context"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"reflect"
)

type Result[T any] interface {
	error
	Value() T
}

// Collect reads from the input channel and collects elements into a slice, respecting context cancellation
func Collect[T any](in <-chan T, ctx context.Context) ([]T, error) {
	var result []T
	for {
		select {
		case <-ctx.Done():
			logger.Log.Warnf("Timeout exceeded, returning %s of size %d", reflect.TypeOf(result), len(result))
			return result, ctx.Err()
		case v, ok := <-in:
			if !ok {
				logger.Log.Debugf("Collected %s of size %d from input channel", reflect.TypeOf(result), len(result))
				// Channel closed, return the collected result
				return result, nil
			}
			result = append(result, v)
		}
	}
}

// FanOut fans a input channel out into two channels, respecting context cancellation
func FanOut[T any](in <-chan T, ctx context.Context) (<-chan T, <-chan T) {
	o1 := make(chan T)
	o2 := make(chan T)
	go func() {
		defer close(o1)
		defer close(o2)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				// Non-blocking send using select to prevent goroutine leak
				select {
				case o1 <- v:
				case <-ctx.Done():
					return
				}
				select {
				case o2 <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return o1, o2
}

// FilterError filters errored Results from channel and calls onError for each, respecting context cancellation
func FilterError[T any](resChan <-chan Result[*T], onError func(err error), context context.Context) <-chan T {
	out := make(chan T)
	go func(onError func(err error)) {
		defer close(out)
		for {
			select {
			case <-context.Done():
				return
			case res, ok := <-resChan:
				if !ok {
					return
				}
				out <- *res.Val
			}
		}
	}(onError)

	return out
}
