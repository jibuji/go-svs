package main

import (
	"fmt"

	"github.com/jibuji/go-svs/stream"
)

func main() {
	words := stream.Distinct(
		stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
	).ToSlice()

	fmt.Printf("Deduplicated words: %v\n", words)
}
