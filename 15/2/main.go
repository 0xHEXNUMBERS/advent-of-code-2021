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
	defer func() {
		sort.Sort(*p)
	}()
	for i, tmp := range *p {
		if tmp.p == point {
			(*p)[i].totalRisk = risk
			return
		}
	}
	*p = append(*p, PointRisk{point, risk})
}

func (p *PriorityQueue) Dequeue() PointRisk {
	ret := (*p)[0]
	*p = (*p)[1:]
	return ret
}

//WARNING: VERY SLOW FOR PART 2
//Took me about 10-15 minutes to fully complete, but it works
func aStar(board [][]int) ([]Point, int) {
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

	totalRisk := 0
	for {
		pRisk := q.Dequeue()
		if pRisk.p == endingPoint {
			totalRisk = pRisk.totalRisk
			break
		}

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
				q.SetRisk(newP, alt+(len(board)-newP.x-1)+(len(board[0])-newP.y-1))
			}
		}
	}
	return nil, totalRisk
}

const SIZE_INCREASE = 5

func setIncreasedValues(board [][]int, baseYLen, baseXLen, basei, basej, start int) {
	for i := 0; i < SIZE_INCREASE; i++ {
		for j := 0; j < SIZE_INCREASE; j++ {
			newValue := (((start - 1) + i + j) % 9) + 1
			board[i*baseYLen+basei][j*baseXLen+basej] = newValue
		}
	}
}

func makeCorrectBoard(base [][]int) [][]int {
	board := make([][]int, len(base)*SIZE_INCREASE)
	for i := range board {
		board[i] = make([]int, len(base[0])*SIZE_INCREASE)
	}

	for i, row := range base {
		for j := range row {
			setIncreasedValues(board, len(base), len(base[0]), i, j, base[i][j])
		}
	}
	return board
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

	correctedBoard := makeCorrectBoard(board)

	_, totalRisk := aStar(correctedBoard)
	fmt.Println(totalRisk)
}
