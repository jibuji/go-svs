package stream

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func (s stream[T]) Parallel(p int) stream[T] {
	return stream[T]{
		parallel: max(p, 1),
		nextFn:   s.nextFn,
	}
}

func (s stream[T]) Take(n int) stream[T] {
	return stream[T]{
		parallel: s.parallel,
		nextFn: func() (T, bool) {
			if n <= 0 {
				var zeroVal T
				return zeroVal, false
			}
			n--
			return s.nextFn()
		},
	}
}

// Sequential currently is not that much useful, supress it for now
func (s stream[T]) _Sequential(bufSize int) stream[T] {
	if s.parallel <= 1 {
		return s
	}
	resCh := make(chan T, bufSize)
	countDown := int32(s.parallel)
	doParallel := func() {
		for i := 0; i < s.parallel; i++ {
			go func() {
				defer func() {
					if atomic.AddInt32(&countDown, -1) <= 0 {
						close(resCh)
					}
				}()
				for {
					v, hasNext := s.nextFn()
					if !hasNext {
						return
					}
					resCh <- v
				}
			}()
		}
	}
	// fmt.Println("seq parallel:", s.parallel)
	// debug.PrintStack()
	go doParallel()
	return stream[T]{
		parallel: 1,
		nextFn: func() (T, bool) {
			v, ok := <-resCh
			return v, ok
		},
	}
}

func (s stream[T]) Filter(predicate func(T) bool) stream[T] {
	return stream[T]{
		parallel: s.parallel,
		nextFn: func() (T, bool) {
			for {
				v, hasNext := s.nextFn()
				if !hasNext {
					return v, false
				}
				if predicate(v) {
					return v, true
				}
			}
		},
	}
}

func (s stream[T]) FilterN(n int, predicate func(T) bool) stream[T] {
	var zeroVal T
	return stream[T]{
		parallel: s.parallel,
		nextFn: func() (T, bool) {
			if n <= 0 {
				return zeroVal, false
			}
			n--
			for {
				v, hasNext := s.nextFn()
				if !hasNext {
					return zeroVal, false
				}
				if predicate(v) {
					return v, true
				}
			}
		},
	}
}

func (s stream[T]) Map(fn func(T) T) stream[T] {
	return Map[T, T](s, fn)
}

func Map[I any, O any](s stream[I], fn func(I) O) stream[O] {
	return stream[O]{
		parallel: s.parallel,
		nextFn: func() (O, bool) {
			v, hasNext := s.nextFn()
			var zeroVal O
			if !hasNext {
				return zeroVal, false
			}
			return fn(v), true
		},
	}
}

// Limit returns a stream consisting of the elements of this stream, truncated to
// be no longer than maxSize in length.
// This function is equivalent to invoking input.Limit(maxSize) as method.
func Limit[T any](s stream[T], n int) stream[T] {
	return s.Take(n)
}

func (s stream[T]) Limit(maxSize int) stream[T] {
	return s.Take(maxSize)
}

func Distinct[T comparable](s stream[T]) stream[T] {
	register := make(map[T]struct{})
	return stream[T]{
		parallel: s.parallel,
		nextFn: func() (T, bool) {
			for {
				v, hasNext := s.nextFn()
				if !hasNext {
					return v, false
				}
				if _, ok := register[v]; !ok {
					register[v] = struct{}{}
					return v, true
				}
			}
		},
	}
}

const ParallelMergeSortThreshold = 1 << 12

// parallelInPlaceMergeSort sorts the given slice in-place using a parallel merge sort algorithm.
// The algorithm uses the specified number of parallel goroutines to perform the merge sort.
func parallelMergeSort[T any](slice []T, comparator Comparator[T], parallel int) {
	if len(slice) <= ParallelMergeSortThreshold || parallel == 1 {
		SortSlice(slice, comparator)
		return
	}

	// Split the slice into two halves.
	mid := len(slice) / 2
	left := slice[:mid]
	right := slice[mid:]

	// Sort the two halves in parallel.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		parallelMergeSort(left, comparator, parallel/2)
	}()

	parallelMergeSort(right, comparator, parallel/2)
	wg.Wait()
	fmt.Println("before merge:", slice, mid)
	// Merge the two sorted halves in-place.
	mergeSortedHalves(slice, mid, comparator)
}

// // [1,2,3,4,5,1,2,3,4,5]
// // merge two sorted parts into one sorted slice in place, the two parts are [0, mid) and [mid, len(slice))
// mid is next of the last element of the first part
func mergeSortedHalves[T any](slice []T, mid int, comparator Comparator[T]) {
	if comparator(slice[mid-1], slice[mid]) <= 0 {
		return
	}
	i, j, high := 0, mid, len(slice)
	for i < mid && j < high {
		if comparator(slice[i], slice[j]) <= 0 {
			i++
		} else {
			tmp := slice[j]
			copy(slice[i+1:], slice[i:j])
			slice[i] = tmp
			i++
			j++
			mid++
		}
	}
}

// Sorted returns a stream consisting of the elements of this stream, sorted according
// to the provided Comparator.
// This function is equivalent to invoking s.Sorted(comparator) as method.
func Sorted[T any](s stream[T], comparator Comparator[T]) stream[T] {
	return s.Sorted(comparator)
}

func (s stream[T]) Sorted(comparator Comparator[T]) stream[T] {
	var elems []T
	doSort := func() {
		elems = s.ToSlice()
		if s.parallel > 1 {
			parallelMergeSort(elems, comparator, s.parallel)
		} else {
			SortSlice(elems, comparator)
		}
		// fmt.Println("sorted:", elems)
	}
	once := sync.Once{}
	index := int64(0)
	return stream[T]{
		// sorted now, we should not parallel afterwards
		// execept user force to do so
		nextFn: func() (T, bool) {
			once.Do(doSort)
			index := atomic.AddInt64(&index, 1)
			// fmt.Println("index:", index, elems)
			if index > int64(len(elems)) {
				var zeroVal T
				return zeroVal, false
			}
			v := elems[index-1]
			return v, true
		},
	}
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to
// each element. Each mapped stream is closed after its contents have been placed into this
// stream. (If a mapped stream is null an empty stream is used, instead.)
//
// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
//
// When both the input and output type are the same, the operation can be
// invoked as the method input.FlatMap(mapper).
func FlatMap[IN, OUT any](input stream[IN], mapper func(IN) stream[OUT]) stream[OUT] {
	nextFromInputStream := input.nextFn
	var nextFromOutputStream func() (OUT, bool)
	return stream[OUT]{
		parallel: input.parallel,
		nextFn: func() (OUT, bool) {
			for {
				if nextFromOutputStream == nil {
					v, hasNext := nextFromInputStream()
					if !hasNext {
						var zeroVal OUT
						return zeroVal, false
					}
					nextFromOutputStream = mapper(v).nextFn
				}
				v, hasNext := nextFromOutputStream()
				if hasNext {
					return v, true
				}
				nextFromOutputStream = nil
			}
		},
	}
}

func (s stream[T]) FlatMap(mapper func(T) stream[T]) stream[T] {
	return FlatMap[T, T](s, mapper)
}

// Peek peturns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
// This function is equivalent to invoking input.Peek(consumer) as method.
// For parallel stream pipelines,
// the action may be called at whatever time and in whatever thread the element is made available by the upstream operation. If the action modifies shared state, it is responsible for providing the required synchronization.
func Peek[T any](input stream[T], consumer func(T)) stream[T] {
	return input.Peek(consumer)
}
func (s stream[T]) Peek(consumer func(T)) stream[T] {
	return stream[T]{
		parallel: s.parallel,
		nextFn: func() (T, bool) {
			v, hasNext := s.nextFn()
			if hasNext {
				consumer(v)
			}
			return v, hasNext
		},
	}
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding
// the first n elements of the stream.
// This function is equivalent to invoking input.Skip(n) as method.
func Skip[T any](input stream[T], n int) stream[T] {
	return input.Skip(n)
}

func (s stream[T]) Skip(n int) stream[T] {
	return stream[T]{
		nextFn: func() (T, bool) {
			// n is a closure here
			for ; n > 0; n-- {
				_, hasNext := s.nextFn()
				if !hasNext {
					var zeroVal T
					return zeroVal, false
				}
			}
			return s.nextFn()
		},
	}
}
