package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	y int
	x int
}

type Direction int

const (
	YDIR Direction = iota
	XDIR
)

type Fold struct {
	dir  Direction
	line int
}

func foldPaper(paper map[Point]struct{}, fold Fold) map[Point]struct{} {
	foldedPaper := make(map[Point]struct{})
	for point := range paper {
		newPoint := point
		if fold.dir == YDIR && point.y > fold.line {
			newPoint.y -= (point.y - fold.line) * 2
		} else if fold.dir == XDIR && point.x > fold.line {
			newPoint.x -= (point.x - fold.line) * 2
		}
		foldedPaper[newPoint] = struct{}{}
	}
	return foldedPaper
}

func paperString(paper map[Point]struct{}) string {
	maxY := 0
	maxX := 0
	for point := range paper {
		if point.y > maxY {
			maxY = point.y
		}
		if point.x > maxX {
			maxX = point.x
		}
	}
	board := make([][]byte, maxY+1)
	for i := range board {
		board[i] = make([]byte, maxX+1)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	for point := range paper {
		board[point.y][point.x] = '#'
	}

	str := ""
	for i := range board {
		str += string(board[i]) + "\n"
	}
	return str
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	sections := strings.Split(string(input), "\n\n")
	coordsData := sections[0]
	foldsData := sections[1]

	paper := make(map[Point]struct{})
	for _, coordData := range strings.Split(coordsData, "\n") {
		xy := strings.Split(coordData, ",")
		x, _ := strconv.Atoi(xy[0])
		y, _ := strconv.Atoi(xy[1])
		paper[Point{y, x}] = struct{}{}
	}

	for _, foldData := range strings.Split(foldsData, "\n") {
		var dirStr byte
		var dir Direction
		var line int
		fmt.Sscanf(foldData, "fold along %c=%d", &dirStr, &line)

		switch dirStr {
		case 'y':
			dir = YDIR
		case 'x':
			dir = XDIR
		}
		paper = foldPaper(paper, Fold{dir, line})
	}

	fmt.Println(paperString(paper))
}
