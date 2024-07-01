// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"fmt"
	"reflect"
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

	received := <-mat
	ptrele.Value = 22

	empty := element{}
	if received == empty {
		t.Errorf("element was empty")
	}

	if fmt.Sprintf("%p", &received) == fmt.Sprintf("%p", &ptrele) {
		t.Errorf("expected different addresses, got %p and %p", &received, &ptrele)
	}

	if received.Value == ptrele.Value {
		t.Errorf("the value was different, got %d and %d", received.Value, ptrele.Value)
	}

}

func TestFanOut(t *testing.T) {
	l := 10
	ch := make(chan int)

	go func() {
		for i := 0; i < l; i++ {
			ch <- 1
		}
		close(ch)
	}()

	c1, c2 := FanOut(ch, context.TODO())

	col1, _ := Collect(c1, context.TODO())
	col2, _ := Collect(c2, context.TODO())

	if len(col1) != len(col2) {
		t.Errorf("got %d, want %d", len(col1), len(col2))
	}
	b := reflect.DeepEqual(col1, col2)

	if !b {
		t.Errorf("got %p, want %p", col1, col2)
	}

}

func TestCollect(t *testing.T) {
	l := 10
	ch := make(chan int, l)

	for i := 0; i < l; i++ {
		ch <- 1
	}
	close(ch)

	col, _ := Collect(ch, context.TODO())

	if len(col) != l {
		t.Errorf("got %d, want %d", len(col), l)
	}
}
