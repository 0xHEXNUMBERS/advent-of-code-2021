package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CompareStrategy int

const (
	MostCommon CompareStrategy = iota
	LeastCommon
)

func getMatch(lines []string, i int, c CompareStrategy) int {
	if len(lines) == 1 {
		ret, _ := strconv.ParseInt(lines[0], 2, 64)
		return int(ret)
	}

	var onesCount int
	for _, line := range lines {
		if line[i] == '1' {
			onesCount++
		}
	}

	var cmp byte
	if float64(onesCount) >= float64(len(lines))/2.0 {
		switch c {
		case MostCommon:
			cmp = '1'
		case LeastCommon:
			cmp = '0'
		}
	} else {
		switch c {
		case MostCommon:
			cmp = '0'
		case LeastCommon:
			cmp = '1'
		}
	}

	newLines := make([]string, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}

		if line[i] == cmp {
			newLines = append(newLines, line)
		}
	}
	return getMatch(newLines, i+1, c)
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	o2 := getMatch(lines, 0, MostCommon)
	co2 := getMatch(lines, 0, LeastCommon)

	fmt.Println(o2 * co2)
}
