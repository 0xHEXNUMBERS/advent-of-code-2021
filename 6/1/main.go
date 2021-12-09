package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func makeFishSlice(fishStr []string) []int {
	fishes := make([]int, len(fishStr))
	for i, n := range fishStr {
		x, _ := strconv.Atoi(n)
		fishes[i] = x
	}
	return fishes
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	input = input[:len(input)-1] //Remove \n
	if err != nil {
		panic(err)
	}

	fishStr := strings.Split(string(input), ",")
	fishes := makeFishSlice(fishStr)

	for round := 0; round < 80; round++ {
		newFishes := make([]int, 0)
		for _, fish := range fishes {
			if fish == 0 {
				newFishes = append(newFishes, 6, 8)
			} else {
				newFishes = append(newFishes, fish-1)
			}
		}
		fishes = newFishes
	}
	fmt.Println(len(fishes))
}
