package pool

type Task[T any] func() (T, error)
