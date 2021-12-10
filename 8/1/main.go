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
	input = input[:len(input)-1] //Remove last \n

	lines := strings.Split(string(input), "\n")
	sum := 0
	for _, line := range lines {
		args := strings.Split(line, " | ")
		//inputSegments := strings.Split(args[0], " ")
		outputSegments := strings.Split(args[1], " ")
		for _, segments := range outputSegments {
			switch len(segments) {
			case 2, 3, 4, 7:
				sum++
			}
		}
	}
	fmt.Println(sum)
}
