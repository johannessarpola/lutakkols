// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
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

func TestFanoIn(t *testing.T) {
	maxl := 10
	cs := []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
	}

	for i, ch := range cs {
		go func(i int, ch chan int) {
			defer close(ch)
			for j := 0; j < maxl; j++ {
				ch <- j * ((i + 1) * 10)
			}
		}(i, ch)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	combined := FanIn[int](ctx, cs...)

	var res []int
	for nbr := range combined {
		res = append(res, nbr)
	}

	if len(res) == 0 {
		t.Errorf("got empty results")
	}

	if len(res) != maxl*len(cs) {
		t.Errorf("got %d, want %d", len(res), maxl*len(cs))
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

func TestFilter(t *testing.T) {
	l := 10
	errs := 2
	ch := make(chan Result[int])

	go func() {
		ch <- Result[int]{
			Err: errors.New("error 1"),
		}

		for i := 0; i < l; i++ {
			ch <- Result[int]{
				Val: 1,
				Err: nil,
			}
		}

		ch <- Result[int]{
			Err: errors.New("error 2"),
		}
		defer close(ch)
	}()

	calls := 0
	f := FilterError(ch, func(err error) {
		calls += 1
	}, context.TODO())

	c, _ := Collect(f, context.TODO())
	if len(c) != l {
		t.Errorf("got %d, want %d", len(c), l)
	}

	if calls != 2 {
		t.Errorf("error callback should been called %d times got %d", errs, calls)
	}

}

func TestMap(t *testing.T) {
	l := 10
	ch := make(chan int)

	go func() {

		for i := 0; i < l; i++ {
			ch <- i
		}

		defer close(ch)
	}()

	mapped := Map(ch, func(n int) (string, error) {
		return fmt.Sprintf("nbr - %d", n), nil
	}, context.TODO())

	var vals []string
	var errs []error
	for s := range mapped {
		vals = append(vals, s.Val)
		if s.Err != nil {
			errs = append(errs, s.Err)
		}
	}

	if len(vals) != l {
		t.Errorf("got %d values, want %d", len(vals), l)
	}

	if len(errs) != 0 {
		t.Errorf("expected no errors got %d", len(errs))
	}
}
