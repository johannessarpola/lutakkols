// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"sync"
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
	}, context.TODO())
}

type element struct {
	Value int
}

func TestMaterialize(t *testing.T) {
	ch := make(chan *element, 1)
	ptrele := element{Value: 10}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {

		ch <- &ptrele
		close(ch)
		wg.Done()
	}()
	wg.Wait()

	mat := Materialize(ch, context.TODO())

	received := element{}
	for ele := range mat {
		received = ele
		break

	}
	ptrele.Value = 22

	if received.Value == ptrele.Value {
		t.Errorf("the value was different, got %d and %d", received.Value, ptrele.Value)
	}

}
