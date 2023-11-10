package stream

import (
	"sync"
	"sync/atomic"

	"golang.org/x/exp/constraints"
)

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

func Generate[T any](fn func() (T, bool)) stream[T] {
	return stream[T]{
		parallel: 1,
		nextFn:   fn,
	}
}

// Range returns a stream of integers from start (inclusive) to end (exclusive)
// The generating process is thread-safe
func Range[T constraints.Integer](start, end T) stream[T] {
	lock := &sync.Mutex{}
	return stream[T]{
		parallel: 1,
		nextFn: func() (T, bool) {
			lock.Lock()
			defer lock.Unlock()
			if start >= end {
				return 0, false
			}
			result := start
			start++
			return result, true
		},
	}
}
