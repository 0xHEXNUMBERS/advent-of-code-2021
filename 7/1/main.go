package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func makeCrabSlice(crabStr []string) []int {
	crabs := make([]int, len(crabStr))
	for i, n := range crabStr {
		x, _ := strconv.Atoi(n)
		crabs[i] = x
	}
	return crabs
}

func minmax(ints []int) (int, int) {
	min, max := ints[0], ints[0]
	for _, i := range ints {
		if min > i {
			min = i
		}
		if max < i {
			max = i
		}
	}
	return min, max
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func calculateCost(crabPos []int, target int) int {
	sum := 0
	for _, pos := range crabPos {
		sum += abs(pos - target)
	}
	return sum
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //remove \n

	crabStr := strings.Split(string(input), ",")
	crabPos := makeCrabSlice(crabStr)

	min, max := minmax(crabPos)

	minCost := calculateCost(crabPos, min)
	for i := min + 1; i <= max; i++ {
		cost := calculateCost(crabPos, i)
		if minCost > cost {
			minCost = cost
		}
	}
	fmt.Println(minCost)
}
