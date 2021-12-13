package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//idk why this works.
func findCorruptedRune(in string) rune {
	tmp := in
	for {
		old := tmp
		tmp = strings.ReplaceAll(tmp, "()", "")
		tmp = strings.ReplaceAll(tmp, "{}", "")
		tmp = strings.ReplaceAll(tmp, "[]", "")
		tmp = strings.ReplaceAll(tmp, "<>", "")

		if tmp == old {
			break
		}
	}
	for _, c := range tmp {
		switch c {
		case ')', '}', ']', '>':
			return c
		}
	}
	return 0
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	sum := 0
	for _, line := range strings.Split(string(input), "\n") {
		switch findCorruptedRune(line) {
		case ')':
			sum += 3
		case ']':
			sum += 57
		case '}':
			sum += 1197
		case '>':
			sum += 25137
		}
	}
	fmt.Println(sum)
}
