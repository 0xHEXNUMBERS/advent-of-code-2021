package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	numbers [5][5]string
	picked  [5][5]bool
}

func (b Board) String() string {
	str := ""
	for _, row := range b.numbers {
		for _, num := range row {
			str += "|" + num
		}
		str += "|\n"
	}
	return str
}

func (b Board) isRowCol(i, j int) bool {
	return (b.picked[0][j] && b.picked[1][j] && b.picked[2][j] &&
		b.picked[3][j] && b.picked[4][j]) ||
		(b.picked[i][0] && b.picked[i][1] && b.picked[i][2] &&
			b.picked[i][3] && b.picked[i][4])
}

func (b Board) sumOfUnmarkedNubmers() int {
	sum := 0
	for i, row := range b.numbers {
		for j, num := range row {
			if !b.picked[i][j] {
				n, _ := strconv.Atoi(num)
				sum += n
			}
		}
	}
	return sum
}

func (b *Board) SetNumber(n string) bool {
	for i, row := range b.numbers {
		for j, num := range row {
			if num == n {
				b.picked[i][j] = true
				return b.isRowCol(i, j)
			}
		}
	}
	return false
}

func BoardFromString(data string) Board {
	var board Board

	for i, line := range strings.Split(data, "\n") {
		j := 0
		for _, n := range strings.Split(line, " ") {
			if n == "" {
				continue
			}
			board.numbers[i][j] = n
			j++
		}
	}
	return board
}

func main() {
	lines, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines = lines[:len(lines)-1] //remove ending \n

	args := strings.Split(string(lines), "\n\n")

	numbersDrawn := strings.Split(args[0], ",")
	boards := make([]Board, len(args)-1)
	for i, boardData := range args[1:] {
		boards[i] = BoardFromString(boardData)
	}

	for _, n := range numbersDrawn {
		for i := range boards {
			if boards[i].SetNumber(n) {
				num, _ := strconv.Atoi(n)
				fmt.Println(num * boards[i].sumOfUnmarkedNubmers())
				return
			}
		}
	}
}
