package main

import (
	"fmt"

	"github.com/jibuji/go-svs/stream"
)

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
