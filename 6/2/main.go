package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	FISH_MAX_COUNTER = 9
	FISH_RESTART     = 6
	TOTAL_DAYS       = 256
)

func makeFishSlice(fishStr []string) []int {
	fishes := make([]int, len(fishStr))
	for i, n := range fishStr {
		x, _ := strconv.Atoi(n)
		fishes[i] = x
	}
	return fishes
}

func initializeFishCounterCounts(fishCounters []int) []int {
	counts := make([]int, FISH_MAX_COUNTER)

	for _, fish := range fishCounters {
		counts[fish]++
	}
	return counts
}

func sum(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	input = input[:len(input)-1] //Remove \n
	if err != nil {
		panic(err)
	}

	fishStr := strings.Split(string(input), ",")
	fishes := makeFishSlice(fishStr)
	counts := initializeFishCounterCounts(fishes)

	for i := 0; i < TOTAL_DAYS; i++ {
		tmp8s := counts[8]
		tmp0s := counts[0]

		counts[8] = counts[0]
		counts[0] = counts[1]
		counts[1] = counts[2]
		counts[2] = counts[3]
		counts[3] = counts[4]
		counts[4] = counts[5]
		counts[5] = counts[6]
		counts[6] = counts[7] + tmp0s
		counts[7] = tmp8s
		counts[8] = tmp0s
	}
	fmt.Println(sum(counts))
}
