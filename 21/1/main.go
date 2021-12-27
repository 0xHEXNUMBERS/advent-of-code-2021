package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Dice interface {
	Roll() int
}

type DeterministicDice struct {
	i int
}

func (d *DeterministicDice) Roll() int {
	d.i = (d.i % 100) + 1
	return d.i
}

type Game struct {
	pos       [2]int
	scores    [2]int
	dice      Dice
	diceRolls int
}

func (g *Game) performTurn(player int) bool {
	g.diceRolls += 3
	roll1 := g.dice.Roll()
	roll2 := g.dice.Roll()
	roll3 := g.dice.Roll()
	moves := roll1 + roll2 + roll3
	newPos := (((g.pos[player] - 1) + moves) % 10) + 1
	g.pos[player] = newPos
	g.scores[player] += newPos
	return g.scores[player] >= 1000
}

func playGame(p0, p1 int, d Dice) (Game, int) {
	g := Game{pos: [2]int{p0, p1}, dice: d}
	for {
		if g.performTurn(0) {
			return g, 0
		} else if g.performTurn(1) {
			return g, 1
		}
	}
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

	stats, winner := playGame(p0Pos, p1Pos, &DeterministicDice{})
	fmt.Println(stats.scores[(winner+1)%2] * stats.diceRolls)
}
