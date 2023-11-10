package main

import (
	"fmt"

	"github.com/jibuji/go-svs/stream"
)

func main() {
	fac8 := stream.Range(1, 9).
		Reduce(1, func(a, b int) int {
			return a * b
		})
	fmt.Println("The factorial of 8 is", fac8)
}
