package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ROUNDS = 10

func executeRound(state string, rules map[string]string) string {
	newState := ""
	for i := range state {
		newState += string(state[i])

		if len(state)-1 == i {
			break
		}

		if res, ok := rules[state[i:i+2]]; ok {
			newState += res
		}
	}
	return newState
}

func minmax(in string) int {
	var min, max int

	charCounts := make(map[rune]int)
	for _, c := range in {
		charCounts[c]++
		min, max = charCounts[c], charCounts[c]
	}

	for _, count := range charCounts {
		if min > count {
			min = count
		}
		if max < count {
			max = count
		}
	}
	return max - min
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	sections := strings.Split(string(input), "\n\n")
	currentState := sections[0]
	rulesData := sections[1]

	rules := make(map[string]string)
	for _, ruleData := range strings.Split(rulesData, "\n") {
		ruleSections := strings.Split(ruleData, " -> ")
		rules[ruleSections[0]] = ruleSections[1]
	}

	for i := 0; i < ROUNDS; i++ {
		currentState = executeRound(currentState, rules)
	}

	fmt.Println(minmax(currentState))
}
