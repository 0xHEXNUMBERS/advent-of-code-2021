package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func inBounds(board [][]int, i, j int) bool {
	return i >= 0 && j >= 0 && i < len(board) && j < len(board[0])
}

func isLowPoint(board [][]int, i, j int) bool {
	diff := [][]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, d := range diff {
		if inBounds(board, i+d[0], j+d[1]) &&
			board[i][j] >= board[i+d[0]][j+d[1]] {
			return false
		}
	}
	return true
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last \n

	rows := strings.Split(string(input), "\n")
	board := make([][]int, len(rows))
	for i, row := range rows {
		board[i] = make([]int, len(row))
		for j, char := range row {
			board[i][j] = int(char - int32('0'))
		}
	}

	sum := 0
	for i, row := range rows {
		for j := range row {
			if isLowPoint(board, i, j) {
				sum += 1 + board[i][j]
			}
		}
	}
	fmt.Println(sum)
}
