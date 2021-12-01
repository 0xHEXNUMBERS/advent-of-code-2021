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

	depth1, _ := strconv.Atoi(lines[0])
	depth2, _ := strconv.Atoi(lines[1])
	depth3, _ := strconv.Atoi(lines[2])
	totalDepth := depth1 + depth2 + depth3

	timesIncreased := 0
	for _, line := range lines[3:] {
		curDepth, _ := strconv.Atoi(line)
		if depth2+depth3+curDepth > totalDepth {
			timesIncreased++
		}
		totalDepth = (totalDepth - depth1) + curDepth
		depth1 = depth2
		depth2 = depth3
		depth3 = curDepth
	}
	fmt.Println(timesIncreased)
}
