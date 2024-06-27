package pool

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type ReturnVal struct {
	n int
}

func failTask() (*ReturnVal, error) {
	return &ReturnVal{rand.Int()}, fmt.Errorf("fail")
}

func normalTask() (*ReturnVal, error) {
	return &ReturnVal{rand.Int()}, nil
}

func longTask() (*ReturnVal, error) {
	time.Sleep(time.Millisecond * 100)
	return &ReturnVal{rand.Int()}, nil
}

func TestThreadpoolTimeout(t *testing.T) {
	s := time.Now()
	size := 100

	var jobs []Task[*ReturnVal]
	for i := 0; i < size; i++ {
		jobs = append(jobs, longTask)
	}
	pool := NewWorkSet[*ReturnVal](jobs, 10, time.Millisecond*1)
	results := pool.Collect()
	if len(results) != size {
		t.Errorf("expecting %d resultQueue, got %d", size, len(results))
	}
	for _, r := range results {
		if r.Value != nil {
			t.Errorf("there was val")
		}
		if r.Error == nil {
			t.Errorf("error was nil")
		}
	}

	if time.Now().Sub(s).Milliseconds() >= 1000 {
		t.Errorf("taskQueue didnt get canceled")
	}

}

func TestThreadpool2(t *testing.T) {
	size := 120

	failCounter := 0
	successCounter := 0
	types := []Task[*ReturnVal]{longTask, normalTask, failTask}
	var jobs []Task[*ReturnVal]
	for i := 0; i < size; i++ {
		jobs = append(jobs, types[i%3])
	}
	pool := NewWorkSet[*ReturnVal](jobs, 10, time.Millisecond*1)
	results := pool.Collect()
	if len(results) != size {
		t.Errorf("expecting %d resultQueue, got %d", size, len(results))
	}
	for _, r := range results {
		if r.Value != nil {
			successCounter += 1
		}
		if r.Error != nil {
			failCounter++
		}
	}

	if successCounter != size/3 {
		t.Errorf("invalid number of successe: %d", successCounter)
	}
	if failCounter != size/3*2 {
		t.Errorf("invalid number of fail: %d", failCounter)
	}
}

func TestCollect(t *testing.T) {
	size := 1000

	var jobs []Task[*ReturnVal]
	for i := 0; i < size; i++ {
		jobs = append(jobs, normalTask)
	}
	pool := NewWorkSet[*ReturnVal](jobs, 10, time.Millisecond*1)

	results := pool.Collect()
	log.Printf("resultQueue count: %d", len(results))
	if len(results) != size {
		t.Errorf("invalid number of resultQueue: %d", len(results))
	}
}
