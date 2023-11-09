package stream

import (
	"sync"
	"sync/atomic"
)

// --- terminals
// ForEach invokes the consumer function for each item of the stream.
// This function is equivalent to invoking input.ForEach(consumer) as method.
// The behavior of this operation is explicitly nondeterministic. For parallel stream pipelines,
// this operation does not guarantee to respect the encounter order of the stream,
// as doing so would sacrifice the benefit of parallelism. For any given element,
// the action may be performed at whatever time and in whatever thread the library chooses.
// If the action accesses shared state, it is responsible for providing the required synchronization.
func ForEach[T any](input stream[T], consumer func(T)) {
	input.ForEach(consumer)
}

func (s stream[T]) ForEach(consumer func(T)) {
	if s.parallel <= 1 {
		next := s.nextFn
		for in, ok := next(); ok; in, ok = next() {
			consumer(in)
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(s.parallel)
	for i := 0; i < s.parallel; i++ {
		go func() {
			defer wg.Done()
			next := s.nextFn
			for in, ok := next(); ok; in, ok = next() {
				consumer(in)
			}
		}()
	}
	wg.Wait()
}

// terminal operation
func (s stream[T]) ToSlice() []T {
	//quick path for sequential stream
	if s.parallel <= 1 {
		res := []T{}
		for {
			v, hasNext := s.nextFn()
			if !hasNext {
				return res
			}
			res = append(res, v)
		}
	}
	var wg sync.WaitGroup
	wg.Add(s.parallel)

	resCh := make(chan T, s.parallel)

	for i := 0; i < s.parallel; i++ {
		go func() {
			defer wg.Done()
			for {
				v, hasNext := s.nextFn()
				if !hasNext {
					return
				}
				resCh <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res := []T{}
	for v := range resCh {
		res = append(res, v)
	}
	return res
}

func (s stream[T]) ReduceSequentially(initial T, fn func(T, T) T) T {
	return ReduceSequentially[T, T](s, initial, fn)
}

// ReduceSequentially Performs a reduction on the elements of this stream, using the provided identity value and an associative accumulation function, and returns the reduced value.
// The identity value must be an identity for the accumulator function. This means that for all t, accumulator.apply(identity, t) is equal to t. The accumulator function must be an associative function.
func ReduceSequentially[I any, O any](s stream[I], identity O, accumulator func(O, I) O) O {
	for {
		v, hasNext := s.nextFn()
		if !hasNext {
			return identity
		}
		identity = accumulator(identity, v)
	}

}

// the identity element is both an initial seed value for the reduction and a default result if there are no input elements. The accumulator function takes a partial result and the next element, and produces a new partial result.
// The combiner function combines two partial results to produce a new partial result.
func Reduce[I any, O any](s stream[I], identity O, accumulator func(O, I) O, combiner func(O, O) O) O {
	//quick path for sequential stream
	if s.parallel <= 1 {
		for {
			v, hasNext := s.nextFn()
			if !hasNext {
				return identity
			}
			identity = accumulator(identity, v)
		}
	}
	var wg sync.WaitGroup
	wg.Add(s.parallel)

	resCh := make(chan O, s.parallel)

	for i := 0; i < s.parallel; i++ {
		go func() {
			defer wg.Done()
			res := identity
			for {
				v, hasNext := s.nextFn()
				if !hasNext {
					resCh <- res
					return
				}
				res = accumulator(res, v)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res := identity
	for v := range resCh {
		res = combiner(res, v)
	}
	return res
}

func (s stream[T]) Reduce(initial T, fn func(T, T) T) T {
	return Reduce[T, T](s, initial, fn, fn)
}

// AllMatch returns whether all elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is false, it stops processing
// the rest of the stream.
// If the stream is empty then true is returned and the predicate is not evaluated.
// This function is equivalent to invoking input.AllMatch(predicate) as method
func AllMatch[T any](input stream[T], predicate func(T) bool) bool {
	return input.AllMatch(predicate)
}

func (s stream[T]) AllMatch(predicate func(T) bool) bool {
	next := s.nextFn
	if s.parallel <= 1 {
		for r, ok := next(); ok; r, ok = next() {
			if !predicate(r) {
				return false
			}
		}
		return true
	}
	var wg sync.WaitGroup
	wg.Add(s.parallel)
	done := int32(0)  // 0: not done, 1: done; for short circuit
	match := int32(1) // 0: not match, 1: match
	for i := 0; i < s.parallel; i++ {
		go func() {
			defer wg.Done()
			for r, ok := next(); ok; r, ok = next() {
				if atomic.LoadInt32(&done) == 1 {
					// short circuit
					return
				}
				if !predicate(r) {
					atomic.StoreInt32(&match, 0)
					atomic.StoreInt32(&done, 1)
					return
				}
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt32(&match) == 1
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
// If the stream is empty then false is returned and the predicate is not evaluated.
// This function is equivalent to invoking input.AnyMatch(predicate) as method
func AnyMatch[T any](input stream[T], predicate func(T) bool) bool {
	return input.AnyMatch(predicate)
}

func (s stream[T]) AnyMatch(predicate func(T) bool) bool {
	next := s.nextFn
	if s.parallel <= 1 {
		for r, ok := next(); ok; r, ok = next() {
			if predicate(r) {
				return true
			}
		}
		return false
	}
	var wg sync.WaitGroup
	wg.Add(s.parallel)
	done := int32(0)  // 0: not done, 1: done; for short circuit
	match := int32(0) // 0: not match, 1: match
	for i := 0; i < s.parallel; i++ {
		go func() {
			defer wg.Done()
			for r, ok := next(); ok; r, ok = next() {
				if atomic.LoadInt32(&done) == 1 {
					// short circuit
					return
				}
				if predicate(r) {
					atomic.StoreInt32(&match, 1)
					atomic.StoreInt32(&done, 1)
					return
				}
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt32(&match) == 1
}

// NoneMatch returns whether no elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
// This function is equivalent to invoking input.NoneMatch(predicate) as method.
func NoneMatch[T any](input stream[T], predicate func(T) bool) bool {
	return input.NoneMatch(predicate)
}

func (is stream[T]) NoneMatch(predicate func(T) bool) bool {
	return !is.AnyMatch(predicate)
}

// Count of elements in this stream.
func Count[T any](input stream[T]) int {
	return input.Count()
}

// the count of elements in this stream
func (s stream[T]) Count() int {
	// use Reduce to count
	return Reduce[T, int](s, 0,
		func(acc int, in T) int { return acc + 1 },
		func(i1, i2 int) int {
			return i1 + i2
		})
}

// Concat creates a lazily concatenated stream whose elements are all the elements of the first stream followed by all the elements of the second stream.
// The resulting stream is ordered if both of the input streams are ordered, and parallel if either of the input streams is parallel.
// When the resulting stream is closed, the close handlers for both input streams are invoked.
func Concat[T any](s1, s2 stream[T]) stream[T] {
	return stream[T]{
		parallel: max(s1.parallel, s2.parallel),
		nextFn: func() (T, bool) {
			v, hasNext := s1.nextFn()
			if hasNext {
				return v, true
			}
			return s2.nextFn()
		},
	}
}

// Max returns the maximum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
// This function is equivalent to invoking input.Max(cmp) as method.
func Max[T any](input stream[T], cmp Comparator[T]) (T, bool) {
	return input.Max(cmp)
}

func (s stream[T]) Max(cmp Comparator[T]) (T, bool) {
	next := s.nextFn
	var max T
	for n, ok := next(); ok; n, ok = next() {
		if cmp(n, max) > 0 {
			max = n
		}
	}
	return max, true
}

// Min returns the minimum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
// This function is equivalent to invoking input.Min(cmp) as method.
func Min[T any](input stream[T], cmp Comparator[T]) (T, bool) {
	return input.Min(cmp)
}

func (s stream[T]) Min(cmp Comparator[T]) (T, bool) {
	next := s.nextFn
	var min T
	for n, ok := next(); ok; n, ok = next() {
		if cmp(n, min) < 0 {
			min = n
		}
	}
	return min, true
}
