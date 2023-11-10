package main

import (
	"fmt"

	"github.com/jibuji/go-svs/stream"
)

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
