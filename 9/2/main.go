package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Point struct {
	y int
	x int
}

func inBounds(board [][]int, p Point) bool {
	return p.y >= 0 && p.x >= 0 && p.y < len(board) && p.x < len(board[0])
}

func isLowPoint(board [][]int, p Point) bool {
	diff := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, d := range diff {
		newP := Point{p.y + d.y, p.x + d.x}
		if inBounds(board, newP) &&
			board[p.y][p.x] >= board[newP.y][newP.x] {
			return false
		}
	}
	return true
}

func findLowPoints(board [][]int) []Point {
	points := make([]Point, 0)
	for i, row := range board {
		for j := range row {
			p := Point{i, j}
			if isLowPoint(board, p) {
				points = append(points, p)
			}
		}
	}
	return points
}

func sizeOfBasinSearch(board [][]int, p Point, pointsFound map[Point]struct{}) {
	pointsFound[p] = struct{}{}

	diff := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, d := range diff {
		newP := Point{p.y + d.y, p.x + d.x}
		if inBounds(board, newP) && board[newP.y][newP.x] != 9 {
			if _, ok := pointsFound[newP]; !ok {
				sizeOfBasinSearch(board, newP, pointsFound)
			}
		}
	}
}

func sizeOfBasin(board [][]int, p Point) int {
	countedPoints := make(map[Point]struct{})
	sizeOfBasinSearch(board, p, countedPoints)
	return len(countedPoints)
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

	lowPoints := findLowPoints(board)
	basinSizes := make([]int, len(lowPoints))
	for i, p := range lowPoints {
		basinSizes[i] = sizeOfBasin(board, p)
	}

	sort.Ints(basinSizes)
	fmt.Println(basinSizes[len(basinSizes)-1] *
		basinSizes[len(basinSizes)-2] *
		basinSizes[len(basinSizes)-3])
}
