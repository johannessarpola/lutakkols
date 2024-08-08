package pipes

import "errors"

// Ring simple struct to work as round robin source for a array
type Ring[T any] struct {
	elements []T
	next     T
	size     int
	idx      int
}

// NewRing creates a new simple ring buffer struct for array
func NewRing[T any](elements []T) (*Ring[T], error) {
	if len(elements) == 0 {
		return nil, errors.New("array is empty")
	}
	return &Ring[T]{
		elements: elements,
		next:     elements[0],
		size:     len(elements),
		idx:      0,
	}, nil
}

// Next returns the next element tracked with internal indexes, loops back to beginning when at the end
func (r *Ring[T]) Next() T {
	n := r.next

	if r.idx == r.size-1 {
		r.idx = 0
	} else {
		r.idx++
	}

	r.next = r.elements[r.idx]
	return n
}
