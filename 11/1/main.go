package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	ROUNDS = 100
)

type Point struct {
	y int
	x int
}

func inRange(board [][]int, p Point) bool {
	return p.y >= 0 && p.x >= 0 && p.y < len(board) && p.x < len(board[0])
}

func executePoint(board [][]int, p Point) int {
	board[p.y][p.x]++
	if board[p.y][p.x] == 10 {
		flashes := 1
		diff := []Point{
			{-1, -1}, {-1, 0}, {-1, 1}, {0, -1},
			{0, 1}, {1, -1}, {1, 0}, {1, 1},
		}
		for _, d := range diff {
			newP := Point{p.y + d.y, p.x + d.x}
			if inRange(board, newP) {
				flashes += executePoint(board, newP)
			}
		}
		return flashes
	}
	return 0
}

func execute(board [][]int) int {
	flashes := 0
	for i, row := range board {
		for j := range row {
			flashes += executePoint(board, Point{i, j})
		}
	}
	for i, row := range board {
		for j := range row {
			if board[i][j] > 9 {
				board[i][j] = 0
			}
		}
	}
	return flashes
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	rows := strings.Split(string(input), "\n")
	board := make([][]int, len(rows))
	for i, row := range rows {
		board[i] = make([]int, len(row))
		for j, char := range row {
			board[i][j] = int(int32(char) - '0')
		}
	}

	sum := 0
	for i := 0; i < ROUNDS; i++ {
		sum += execute(board)
	}
	fmt.Println(sum)
}
