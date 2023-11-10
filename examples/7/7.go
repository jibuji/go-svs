package main

import (
	"fmt"
	"time"

	"github.com/jibuji/go-svs/stream"
)

func isPrime(n int) bool {
	for i := 2; i < n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	// we do some complex and nonsense calculation in parallel
	result := stream.Range(1, 1000).
		Parallel(4). // parallelize the process and use 4 goroutines, then the following operations will be executed in parallel (if possible)
		Filter(isPrime).
		Map(func(n int) int {
			// simulate a complex calculation
			time.Sleep(100 * time.Nanosecond)
			return (n + 1) / 2
		}).
		Reduce(0, func(a, b int) int {
			// simulate a complex calculation
			time.Sleep(100 * time.Nanosecond)
			return a + b
		})
	fmt.Println("The nonsense result is", result)
}
