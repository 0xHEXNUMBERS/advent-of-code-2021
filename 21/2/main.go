package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	WinningNumber  = 21
	RollsPerPlayer = 3
)

type Game struct {
	pos           [2]int
	scores        [2]int
	currentPlayer int
}

func calculateUniverses(g Game, moves int) (wins [2]int) {
	newPos := (((g.pos[g.currentPlayer] - 1) + moves) % 10) + 1
	g.pos[g.currentPlayer] = newPos
	g.scores[g.currentPlayer] += newPos
	if g.scores[g.currentPlayer] >= WinningNumber {
		wins[g.currentPlayer]++
		return wins
	}
	g.currentPlayer = (g.currentPlayer + 1) % 2

	return generateUniverses(g)
}

func generateUniverses(g Game) (wins [2]int) {
	threeRolls := calculateUniverses(g, 3)
	fourRolls := calculateUniverses(g, 4)
	fiveRolls := calculateUniverses(g, 5)
	sixRolls := calculateUniverses(g, 6)
	sevenRolls := calculateUniverses(g, 7)
	eightRolls := calculateUniverses(g, 8)
	nineRolls := calculateUniverses(g, 9)

	fourRolls[0] *= 3  //3 ways to roll a 4
	fourRolls[1] *= 3  //3 ways to roll a 4
	fiveRolls[0] *= 6  //6 ways to roll a 5
	fiveRolls[1] *= 6  //6 ways to roll a 5
	sixRolls[0] *= 7   //7 ways to roll a 6
	sixRolls[1] *= 7   //7 ways to roll a 6
	sevenRolls[0] *= 6 //6 ways to roll a 7
	sevenRolls[1] *= 6 //6 ways to roll a 7
	eightRolls[0] *= 3 //3 ways to roll a 8
	eightRolls[1] *= 3 //3 ways to roll a 8

	wins[0] += threeRolls[0] + fourRolls[0] + fiveRolls[0] +
		sixRolls[0] + sevenRolls[0] + eightRolls[0] + nineRolls[0]
	wins[1] += threeRolls[1] + fourRolls[1] + fiveRolls[1] +
		sixRolls[1] + sevenRolls[1] + eightRolls[1] + nineRolls[1]
	return wins
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	players := strings.Split(string(input), "\n")
	p0Pos, _ := strconv.Atoi(players[0][28:]) //Remove start of line
	p1Pos, _ := strconv.Atoi(players[1][28:]) //Remove start of line

	possibleWins := generateUniverses(Game{pos: [2]int{p0Pos, p1Pos}})
	if possibleWins[0] > possibleWins[1] {
		fmt.Println(possibleWins[0])
	} else {
		fmt.Println(possibleWins[1])
	}
}
