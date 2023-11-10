package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jibuji/go-svs/stream"
)

func main() {

	cmp := func(a, b uint32) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	fmt.Println("picking up 5 random numbers from lower to higher")

	stream.Generate(func() (uint32, bool) {
		return r.Uint32(), true
	}).
		Limit(5).
		Sorted(cmp).
		ForEach(func(n uint32) {
			fmt.Println(n)
		})

}
