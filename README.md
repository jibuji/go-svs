# AVA-Stream API

Lightweight and super fast functional stream processing library inspired by [mariomac/gostream](https://github.com/mariomac/gostream) and [RxGo](https://github.com/ReactiveX/RxGo).

## Table of contents

- [Requirements](#requirements)
- [Usage examples](#usage-examples)
- [Limitations](#limitations)
- [Performance](#performance)
- [Completion status](#completion-status)
- [Extra credits](#extra-credits)

## Requirements

- Go 1.18 or higher

This library makes intensive usage of [Type Parameters (generics)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) so it is not compatible with any Go version lower than 1.18.

## Usage examples

For more details about the API, please check the `stream/*_test.go` files.

### Example 1: basic creation, transformation and iteration

```go

func main() {
	stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).
		Filter(isPrime).
		Map(square).
		ForEach(func(n int) {
			fmt.Printf("%d is a square of a prime number\n", n)
		})
}

func isPrime(n int) bool {
	for i := 2; i < n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func square(n int) int {
	return n * n
}
```

Output:

```
1 is a prime number
2 is a prime number
3 is a prime number
5 is a prime number
7 is a prime number
11 is a prime number
```

### Example 2: generation, map, limit and slice conversion

1. Creates an **infinite** stream of random integers (no problem, streams are evaluated lazily!)
2. In this case, if you want to turn the stream into a slice, you need to limit the number of elements
   (the stream can be **infinite**, but slice can't).

```go
rand.Seed(time.Now().UnixMilli())
fmt.Println("let me throw 5 times a dice for you")

// Generate an infinite stream of random numbers
// The generator function of the parameter returns a pair (value, hasMore)
results := stream.Generate(func()(int, bool) {
    return rand.Int(), true
}).
    Map(func(n int) int {
        return n%6 + 1
    }).
    Limit(5).
    ToSlice()

fmt.Printf("results: %v\n", results)
```

Output:

```
let me throw 5 times a dice for you
results: [3 5 2 1 3]
```

### Example 3: Generation from an Range, Map to a different type

1. Generates an stream from in a range of numbers, and then map the number to its `doubled` value.
2. Then `Map` the numbers' stream to a strings' stream. Because, at the moment,
   [go does not allow type parameters in methods](https://github.com/golang/go/issues/49085),
   we need to invoke the `stream.Map` function instead of the `numbers.Map` method
   because the contained type of the output stream (`string`) is different than the type of
   the input stream (`int`).
3. Converts the words stream to a slice and prints it.

```go
func main() {
    numbers := stream.Range(1, 7).Map(func(n int) int {
        return double(n)
    })
    words := stream.Map(numbers, asWord).ToSlice()
    fmt.Println(words)
}

func double(n int) int {
    return 2 * n
}

func asWord(n int) string {
    if n < 10 {
        return []string{"zero", "one", "two", "three", "four", "five",
            "six", "seven", "eight", "nine"}[n]
    } else {
        return "many"
    }
}
```

Output:

```
[one two four eight many many]
```

### Example 4: deduplication of elements

Following example requires to compare the elements of the Stream, so the Stream needs to be
composed by `comparable` elements (this is, accepted by the the `==` and `!=` operators):

1. Instantiate a `Stream` of `comparable` items.
2. Pass it to the `Distinct` method, that will return a copy of the original Stream without
   duplicates
3. Operating as any other stream.

```go
words := stream.Distinct(
  stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
).ToSlice()

fmt.Printf("Deduplicated words: %v\n", words)
```

Output:

```
Deduplicated words: [hello ! ho]
```

### Example 5: sorting from higher to lower

1. Generate a stream of uint32 numbers.
2. Picking up 5 elements.
3. Sorting them by the inverse natural order (from higher to lower)
   - It's **important** to limit the number of elements, avoiding invoking
     `Sorted` over an infinite stream (otherwise it would panic).

```go
fmt.Println("picking up 5 random numbers from higher to lower")
stream.Generate(rand.Uint32).
    Limit(5).
    Sorted(order.Inverse(order.Natural[uint32])).
    ForEach(func(n uint32) {
        fmt.Println(n)
    })
```

Output:

```
picking up 5 random numbers from higher to lower
4039455774
2854263694
2596996162
1879968118
1823804162
```

### Example 6: Reduce and helper functions

1. Generate an infinite incremental Stream (1, 2, 3, 4...) using the `stream.Iterate`
   instantiator and the `item.Increment` helper function.
2. Limit the generated to 8 elements
3. Reduce all the elements multiplying them using the item.Multiply helper function

```go
fac8, _ := stream.Range(1, 9).
    Reduce(item.Multiply[int])
fmt.Println("The factorial of 8 is", fac8)
```

Output:

```
The factorial of 8 is 40320
```

### Example 7: Paralleling process

```go
facOfFac, _ := stream.Range(1, 900000).
    Parallel(4). // parallelize the process and use 4 goroutines, then the following operations will be executed in parallel (if possible)
    Map(func(n int) int {
        // do some heavy processing
        return stream.Range(1, n+1).Reduce(1, (a, b int) int{ return a * b})
    }).
    Reduce(1, (a, b int) int{ return a * b})
fmt.Println("The factorial of 8 is", fac8)
```

## Limitations

Due to the initial limitations of Go generics, the API has the following limitations.
We will work on overcome them as long as new features are added to the Go type parameters
specification.

- You can use `Map` and `FlatMap` as method as long as the output element has the same type of the input.
  If you need to map to a different type, you need to use `stream.Map` or `stream.FlatMap` as functions.
- There is no `Distinct` method. There is only a `stream.Distinct` function.

## Performance

For small streams, the performance of this library is comparable to the performance of [go-stream](https://github.com/mariomac/gostream), but for large streams, the performance of this library is much better.

If you enable parallelism, the performance of this library is even better.

## Completion status

- Stream instantiation functions
  - [x] Comparable
  - [x] Concat
  - [x] Generate
  - [x] Of
  - [x] OfSlice
  - [x] OfChannel
- Stream transformers
  - [x] Distinct
  - [x] Filter
  - [x] FlatMap
  - [x] Limit
  - [x] Map
  - [x] Peek
  - [x] Skip
  - [x] Sorted
  - [ ] GroupBy
  - [ ] Defer
  - [ ] Fork
  - [ ] enable user to early terminate the heavy operations
- Collectors/Terminals
  - [x] ToSlice
  - [x] AllMatch
  - [x] AnyMatch
  - [x] Count
  - [x] ForEach
  - [x] Max
  - [x] Min
  - [x] NoneMatch
  - [x] Reduce
  - [x] ReduceSequentially

## Extra credits

The Stream processing and aggregation functions and docs are heavily inspired by the
[mariomac/gostream](<[mariomac/gostream](https://github.com/mariomac/gostream)>).
