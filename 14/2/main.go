package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ROUNDS = 40

type State struct {
	ruleCounts map[string]int
	charCounts map[byte]int
}

func newState() State {
	return State{
		make(map[string]int),
		make(map[byte]int),
	}
}

func executeRound(state State, rules map[string]string) State {
	newState := newState()
	newState.charCounts = state.charCounts
	for rule, res := range rules {
		pair1 := string(rule[0]) + res
		pair2 := res + string(rule[1])

		newState.ruleCounts[pair1] += state.ruleCounts[rule]
		newState.ruleCounts[pair2] += state.ruleCounts[rule]
		newState.charCounts[res[0]] += state.ruleCounts[rule]
	}
	return newState
}

func minmax(in State) int {
	var min, max int

	for _, count := range in.charCounts {
		min, max = count, count
		break
	}

	for _, count := range in.charCounts {
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
	baseString := sections[0]
	rulesData := sections[1]

	rules := make(map[string]string)
	for _, ruleData := range strings.Split(rulesData, "\n") {
		ruleSections := strings.Split(ruleData, " -> ")
		rules[ruleSections[0]] = ruleSections[1]
	}

	currentState := newState()
	for i := 0; i < len(baseString); i++ {
		currentState.charCounts[baseString[i]]++

		if i < len(baseString)-1 {
			currentState.ruleCounts[baseString[i:i+2]]++
		}
	}

	for i := 0; i < ROUNDS; i++ {
		currentState = executeRound(currentState, rules)
	}

	fmt.Println(minmax(currentState))
}
