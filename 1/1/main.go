package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(stdin), "\n")
	prevDepth, _ := strconv.Atoi(lines[0])
	timesIncreased := 0
	for _, line := range lines[1:] {
		curDepth, _ := strconv.Atoi(line)
		if curDepth > prevDepth {
			timesIncreased++
		}
		prevDepth = curDepth
	}
	fmt.Println(timesIncreased)
}
