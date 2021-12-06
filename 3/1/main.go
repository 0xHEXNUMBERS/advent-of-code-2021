package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	onesCount := make([]int, len(lines[0]))
	for _, line := range lines {
		for i, c := range line {
			if c == '1' {
				onesCount[i]++
			}
		}
	}

	gamma := 0
	epsilon := 0
	for _, count := range onesCount {
		gamma <<= 1
		epsilon <<= 1
		if count > len(lines)/2 {
			gamma |= 1
		} else {
			epsilon |= 1
		}
	}
	fmt.Println(gamma * epsilon)
}
