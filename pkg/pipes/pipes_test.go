// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"testing"
)

func dummySink[T any](data []T) error {
	return nil
}

func TestPour(t *testing.T) {
	l := 10
	ch := make(chan int, 10)
	for i := 0; i < l; i++ {
		ch <- 1
	}
	close(ch)

	cntr := 0

	_ = Pour(ch, func(ints []int) error {
		for _, i := range ints {
			cntr += i
		}
		return nil
	}, context.Background())

	if cntr != l {
		t.Errorf("got %d, want %d", cntr, l)
	}

	cntr = 0
	ch = make(chan int)

	go func() {
		for i := 0; i < l; i++ {
			ch <- 1
		}
		close(ch)
	}()

	_ = Pour(ch, func(ints []int) error {
		for _, i := range ints {
			cntr += i
		}
		return nil
	}, context.Background())
}
