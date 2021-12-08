package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Line struct {
	A Point
	B Point
}

func makePoint(pointStr string) Point {
	coords := strings.Split(pointStr, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	return Point{x, y}
}

func makeLines(inputLines []string) []Line {
	lines := make([]Line, len(inputLines))
	for i, l := range inputLines {
		points := strings.Split(l, " -> ")
		lines[i] = Line{
			makePoint(points[0]),
			makePoint(points[1]),
		}
	}
	return lines
}

func Max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func Min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func Sign(i int) int {
	if i > 0 {
		return 1
	} else if i < 0 {
		return -1
	}
	return 0
}

type Board [][]int

func (b Board) addLine(l Line) {
	p := l.A
	for {
		b[p.x][p.y]++
		if p == l.B {
			return
		}

		p.x += Sign(l.B.x - p.x)
		p.y += Sign(l.B.y - p.y)
	}
}

func (b Board) numberOfOverlappingPoints() int {
	points := 0
	for i, row := range b {
		for j := range row {
			if b[i][j] > 1 {
				points++
			}
		}
	}
	return points
}

func makeBoard(lines []Line) Board {
	maxX := 0
	maxY := 0
	for _, l := range lines {
		if maxX < l.A.x {
			maxX = l.A.x
		}
		if maxX < l.B.x {
			maxX = l.B.x
		}
		if maxY < l.A.y {
			maxY = l.A.y
		}
		if maxY < l.B.y {
			maxY = l.B.y
		}
	}
	b := make(Board, maxX+1)
	for i := range b {
		b[i] = make([]int, maxY+1)
	}
	return b
}

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	inputLines := strings.Split(string(stdin), "\n")
	inputLines = inputLines[:len(inputLines)-1]

	lines := makeLines(inputLines)
	board := makeBoard(lines)
	for _, l := range lines {
		board.addLine(l)
	}

	fmt.Println(board.numberOfOverlappingPoints())
}
