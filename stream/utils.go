package stream

import "math/rand"

func randomSlice(i int) []int {
	s := make([]int, i)
	for j := 0; j < i; j++ {
		//random generate int
		s[j] = randomInt(0, i)
	}
	return s
}

func randomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
