// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"reflect"
	"sync"
)

// Result is a struct to wrap either a element or a error
type Result[T any] struct {
	Val T
	Err error
}

// Pour consumes a channel, collects them into array and calls the sink func with it, respecting context cancellation
func Pour[T any](in <-chan T, sink func([]T) error, ctx context.Context, initial ...T) error {
	collect, err := Collect(in, ctx, initial...)
	if err != nil {
		return err
	}
	return sink(collect)
}

// Map transforms elements in channel to another type
func Map[T any, O any](in <-chan T, fn func(T) (O, error), context context.Context) <-chan Result[O] {
	out := make(chan Result[O])
	go func() {
		defer close(out)
		for {
			select {
			case <-context.Done():
				out <- Result[O]{Err: context.Err()}
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							out <- Result[O]{Err: fmt.Errorf("panic in transformation: %v", r)}
						}
					}()
					mappedValue, err := fn(value)
					select {
					case out <- Result[O]{Val: mappedValue, Err: err}:
					case <-context.Done():
						return
					}
				}()
			}
		}
	}()
	return out
}

// Materialize copies the pointer values into another channel as a concrete type, respecting context cancellation
func Materialize[T any](in <-chan *T, context context.Context) <-chan T {
	out := make(chan T)
	go func(in <-chan *T) {
		defer close(out)
		for {
			select {
			case <-context.Done():
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- *value
			}
		}
	}(in)
	return out
}

// Collect reads from the input channel and collects elements into a slice, respecting context cancellation
func Collect[T any](in <-chan T, ctx context.Context, initial ...T) ([]T, error) {
	result := initial
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

func FanIn[T any](ctx context.Context, chans ...chan T) chan T {
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))

	out := make(chan T)
	for _, ch := range chans {
		go func(ch <-chan T, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
			}()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if ok {
						out <- v
					} else {
						return
					}
				}
			}
		}(ch, wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanOut fans a input channel out into two channels, respecting context cancellation
func FanOut[T any](in <-chan T, ctx context.Context) (<-chan T, <-chan T) {
	o1 := make(chan T)
	o2 := make(chan T)
	go func() {
		var arr []T

	consume:
		for {
			select {
			case <-ctx.Done():
				break consume
			case v, ok := <-in:
				if !ok {
					break consume
				}
				arr = append(arr, v)
			}
		}

		go func(arr []T) {
			defer close(o1)

			for _, e := range arr {
				select {
				case o1 <- e:
				case <-ctx.Done():
					return
				}
			}
		}(arr)

		go func(arr []T) {
			defer close(o2)

			for _, e := range arr {
				select {
				case o2 <- e:
				case <-ctx.Done():
					return
				}
			}
		}(arr)

	}()

	return o1, o2
}

// FilterError filters errored Results from channel and calls onError for each, respecting context cancellation
func FilterError[T any](resChan <-chan Result[T], onError func(err error), context context.Context) <-chan T {
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
				if res.Err != nil {
					onError(res.Err)
				} else {
					out <- res.Val
				}
			}
		}
	}(onError)

	return out
}
