// Package workset contains a simple pool to use to queue work concurrently and then aggregate the resultQueue into a single result
package workset

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// WorkSet struct is used to queue work concurrently
type WorkSet[T any] struct {
	taskQueue   chan Task[T]
	resultQueue chan Result[T]
	taskTimeout time.Duration
	workers     int
	waitGroup   sync.WaitGroup
	stopChan    chan struct{}
}

// NewWorkSet workset which starts a bunch of jobs on n number of workers use Collect() to get the results
func NewWorkSet[T any](jobs []Task[T], workers int, timeout time.Duration) *WorkSet[T] {
	workAmount := len(jobs)
	pool := WorkSet[T]{
		// tasks & results should be 1:1 as even a error is a result
		taskQueue:   make(chan Task[T], workAmount),
		resultQueue: make(chan Result[T], workAmount),
		taskTimeout: timeout,
		workers:     workers,
		stopChan:    make(chan struct{}),
	}
	pool.waitGroup.Add(workAmount)
	pool.startWorkers()
	pool.queue(jobs)
	return &pool
}

// queue adds jobs to task queue
func (p *WorkSet[T]) queue(jobs []Task[T]) {
	for _, job := range jobs {
		p.taskQueue <- job
	}
}

// startWorkers starts the worker goroutines
func (p *WorkSet[T]) startWorkers() {
	for i := 0; i < p.workers; i++ {
		workerId := fmt.Sprintf("Worker-%03d", i+1)
		go p.worker(workerId)
	}
}

func taskHandler[T any](task Task[T], result chan Result[T]) {
	now := time.Now()
	var rs Result[T]
	v, err := task()
	if err != nil {
		rs = Result[T]{Error: err}
	} else {
		rs = Result[T]{
			Duration: time.Since(now),
			Value:    v,
		}
	}
	result <- rs

}

func (p *WorkSet[T]) worker(workerId string) {
	for {
		select {
		case task := <-p.taskQueue:
			// single taskResult channel to gather the task taskResult
			taskResult := make(chan Result[T], 1)

			// asynchronously startWorkers task so that we can see that it does not exceed taskTimeout
			go taskHandler(task, taskResult)

			// Execute but queue not exceed taskTimeout
			var result Result[T]
			select {
			case r := <-taskResult:
				r.WorkerId = workerId
				result = r
			case <-time.After(p.taskTimeout):
				result = Result[T]{WorkerId: workerId, Error: errors.New("taskTimeout exceeded")}
			}

			p.resultQueue <- result

		case <-p.stopChan:
			return
		}
	}
}

func onResult[T any](
	resultSink func(Result[T]),
	resultChan <-chan Result[T]) {
	for result := range resultChan {
		resultSink(result)
	}
}

// Collect all results from the tasks queue
func (p *WorkSet[T]) Collect() []Result[T] {
	var results []Result[T]

	go onResult(func(r Result[T]) {
		results = append(results, r)
		// signal done once the result is gathered
		p.waitGroup.Done()
	}, p.resultQueue)

	p.waitGroup.Wait()
	close(p.stopChan)
	close(p.resultQueue)
	close(p.taskQueue)
	return results
}
