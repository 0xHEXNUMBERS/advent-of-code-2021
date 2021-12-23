package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Sign(i int) int {
	if i > 0 {
		return 1
	} else if i < 0 {
		return -1
	}
	return 0
}

type Simulator struct {
	xVel, yVel int
	xMin, yMin int
	xMax, yMax int
	xPos, yPos int
}

func (s Simulator) Simulate() (int, bool) {
	maxY := 0
	for {
		s.xPos += s.xVel
		s.yPos += s.yVel

		s.xVel -= Sign(s.xVel)
		s.yVel--

		if s.yPos > maxY {
			maxY = s.yPos
		}

		if s.xMin <= s.xPos && s.xPos <= s.xMax &&
			s.yMin <= s.yPos && s.yPos <= s.yMax {
			return maxY, true
		}
		if s.xPos > s.xMax || s.yPos < s.yMin {
			return 0, false
		}
	}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	str := strings.TrimPrefix(string(input), "target area: ")
	args := strings.Split(str, ", ")
	xStrs := strings.Split(args[0], "..")
	xStrs[0] = strings.TrimPrefix(xStrs[0], "x=")
	yStrs := strings.Split(args[1], "..")
	yStrs[0] = strings.TrimPrefix(yStrs[0], "y=")

	xMin, _ := strconv.Atoi(xStrs[0])
	xMax, _ := strconv.Atoi(xStrs[1])
	yMin, _ := strconv.Atoi(yStrs[0])
	yMax, _ := strconv.Atoi(yStrs[1])

	count := 0
	for xVel := -1000; xVel < 1000; xVel++ {
		for yVel := -1000; yVel < 1000; yVel++ {
			s := Simulator{xVel, yVel, xMin, yMin, xMax, yMax, 0, 0}
			if _, ok := s.Simulate(); ok {
				count++
			}
		}
	}
	fmt.Println(count)
}
