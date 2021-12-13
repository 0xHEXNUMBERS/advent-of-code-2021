package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

//idk why this works.
func isCorruptedRune(in string) bool {
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
			return true
		}
	}
	return false
}

func reverseString(in string) string {
	out := []rune(in)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

func autoComplete(in string) string {
	tmp := ""
	for _, c := range in {
		switch c {
		case ')', ']', '}', '>':
			tmp = tmp[:len(tmp)-1]
		default:
			tmp += string(c)
		}
	}
	tmp = reverseString(tmp)

	out := ""
	for _, c := range tmp {
		switch c {
		case '(':
			out += ")"
		case '[':
			out += "]"
		case '{':
			out += "}"
		case '<':
			out += ">"
		}
	}
	return out
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	scores := make([]int, 0)
	for _, line := range strings.Split(string(input), "\n") {
		if isCorruptedRune(line) {
			continue
		}

		appendChars := autoComplete(line)

		acc := 0
		for _, c := range appendChars {
			acc *= 5
			switch c {
			case ')':
				acc += 1
			case ']':
				acc += 2
			case '}':
				acc += 3
			case '>':
				acc += 4
			}
		}
		scores = append(scores, acc)
	}

	sort.Ints(scores)
	fmt.Println(scores[len(scores)/2])
}
