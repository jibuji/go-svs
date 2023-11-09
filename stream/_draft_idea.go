package stream

// example:
// stream.ofChannel(ch).map(t => ({gender: t.gender, name: t.name, class: t.class, score: t.score})).
// 	groupBy(t => t.class).sortBy((a, b) => a < b).map([key, s] => s.parallel(2).sortScore((a, b) => a < b).take(3).map(t => ({name: t.name, score: t.score})).toSlice()).toSlice()

type keyedStream[K comparable, T any] struct {
	key K
	s   stream[T]
}

type StreamGroup[K comparable, T any] chan keyedStream[K, T]

func GroupBy1[K comparable, T any](s stream[T], keyFn func(T) K, streamBufSize int, groupBufferSize int) StreamGroup[K, T] {
	register := make(map[K]chan T, max(groupBufferSize, 1))
	sg := make(StreamGroup[K, T])
	for {
		v, hasNext := s.nextFn()
		if !hasNext {
			break
		}
		key := keyFn(v)
		ch, ok := register[key]
		if !ok {
			ch = make(chan T, max(streamBufSize, 0))
			register[key] = ch
			go func() {
				sg <- keyedStream[K, T]{
					key: key,
					s:   OfChannel[T](ch),
				}
			}()
		}
		ch <- v
	}
	return sg
}

// example: s.GroupBy(t => t.class).SortBy((a, b) => a < b).Map([key, s] => s.Parallel(2).SortScore((a, b) => a < b).Take(3).Map(t => ({name: t.name, score: t.score})).ToSlice()).ToSlice()
func GroupBy[K comparable, T any](s stream[T], keyFn func(T) K, streamBufSize int, groupBufferSize int) StreamGroup[K, T] {
	register := make(map[K]chan T, max(groupBufferSize, 1))
	sg := make(StreamGroup[K, T])
	for {
		v, hasNext := s.nextFn()
		if !hasNext {
			break
		}
		key := keyFn(v)
		ch, ok := register[key]
		if !ok {
			ch = make(chan T, max(streamBufSize, 0))
			register[key] = ch
			go func() {
				sg <- keyedStream[K, T]{
					key: key,
					s:   OfChannel[T](ch),
				}
			}()
		}
		ch <- v
	}
	return sg
}

// type ShardKey[T any] interface {
// 	Key(index int, item T) string
// }

// type MultiShards[T any] []stream[T]

// func (s stream[T]) Shard(key ShardKey[T]) MultiShards[T] {
// 	return []stream[T]{}
// }

// func Parallel[In any, Out any](ss MultiShards[In], pNum int, t func(stream[In]) stream[Out]) MultiShards[Out] {
// 	return MultiShards[Out]{}
// }

// func (o MultiShards[T]) SortByKey(key1, key2 string) MultiShards[T] {
// 	return MultiShards[T]{}
// }

// func genTickets(n int) chan struct{} {
// 	tickets := make(chan struct{}, n)
// 	for i := 0; i < n; i++ {
// 		tickets <- struct{}{}
// 	}
// 	return tickets
// }

/**
Problems remained:
1. If users stop consume the stream, the goroutines for parallel ops will be leaked. How to solve this problem?
2. Do we need to handle error in the stream? like, when error happens, stop the stream and return the error.
TODO:
1. Fork
2. GroupBy
3. Defer
4. Window
6. Zip/Unzip
**/
