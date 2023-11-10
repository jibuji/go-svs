package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jibuji/go-svs/stream"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	fmt.Println("let me throw 5 times a dice for you")

	// Generate an infinite stream of random numbers
	// The generator function of the parameter returns a pair (value, hasMore)
	results := stream.Generate(func() (int, bool) {
		return r.Int(), true
	}).
		Map(func(n int) int {
			return n%6 + 1
		}).
		Limit(5).
		ToSlice()

	fmt.Printf("results: %v\n", results)
}
