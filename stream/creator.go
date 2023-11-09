package stream

import "sync/atomic"

type stream[T any] struct {
	// sorted   bool
	parallel int //will be used in terminal operation or stateful operation
	nextFn   func() (T, bool)
}

func Of[T any](elems ...T) stream[T] {
	return OfSlice(elems)
}

func OfSlice[T any](elems []T) stream[T] {
	i := int64(0)
	return stream[T]{
		parallel: 1,
		nextFn: func() (T, bool) {
			i := atomic.AddInt64(&i, 1)
			var v T
			if i > int64(len(elems)) {
				return v, false
			}
			v = elems[i-1]
			return v, true
		},
	}
}

func OfChannel[T any](ch chan T) stream[T] {
	return stream[T]{
		parallel: 1,
		nextFn: func() (T, bool) {
			v, ok := <-ch
			return v, ok
		},
	}
}
