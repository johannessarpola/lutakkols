package workset

type Task[T any] func() (T, error)
