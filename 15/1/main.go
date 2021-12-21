package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

const LARGE_NUMBER = 100000000

type Point struct {
	y int
	x int
}

func (p Point) InBounds(b [][]int) bool {
	return p.y >= 0 && p.y < len(b) && p.x >= 0 && p.x < len(b[0])
}

func (p Point) Add(p2 Point) Point {
	return Point{p.y + p2.y, p.x + p2.x}
}

type PointRisk struct {
	p         Point
	totalRisk int
}

type PriorityQueue []PointRisk

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].totalRisk < p[j].totalRisk
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *PriorityQueue) SetRisk(point Point, risk int) {
	for i, tmp := range *p {
		if tmp.p == point {
			(*p)[i].totalRisk = risk
			return
		}
	}
	*p = append(*p, PointRisk{point, risk})
}

func dijkstra(board [][]int) ([]Point, int) {
	startPoint := Point{0, 0}
	endingPoint := Point{len(board) - 1, len(board[0]) - 1}
	q := PriorityQueue{}

	distances := make(map[Point]int)
	prevPoint := make(map[Point]Point)

	distances[startPoint] = 0
	for i, row := range board {
		for j := range row {
			p := Point{i, j}
			if p != startPoint {
				distances[p] = LARGE_NUMBER
			}
			q.SetRisk(p, distances[p])
		}
	}
	sort.Sort(q)

	pointDiffs := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for pRisk := q[0]; pRisk.p != endingPoint; pRisk = q[0] {
		point := pRisk.p
		for _, d := range pointDiffs {
			newP := point.Add(d)
			if !newP.InBounds(board) {
				continue
			}

			alt := distances[point] + board[newP.y][newP.x]
			if alt < distances[newP] {
				distances[newP] = alt
				prevPoint[newP] = point
				q.SetRisk(newP, alt)
			}
		}
		q = q[1:]
		sort.Sort(q)
	}
	totalRisk := q[0].totalRisk

	return nil, totalRisk
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	lines := strings.Split(string(input), "\n")
	board := make([][]int, len(lines))
	for i, row := range lines {
		board[i] = make([]int, len(row))
		for j, risk := range row {
			board[i][j], _ = strconv.Atoi(string(risk))
		}
	}

	_, totalRisk := dijkstra(board)
	fmt.Println(totalRisk)
}
