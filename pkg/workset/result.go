package workset

import "time"

type Result[T any] struct {
	Duration time.Duration
	Value    T
	Error    error
	WorkerId string
}
