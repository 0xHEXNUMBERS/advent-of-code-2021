package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const CommonPointsNeeded = 12

type Point struct {
	x int
	y int
	z int
}

func (p Point) Add(p2 Point) Point {
	return Point{p.x + p2.x, p.y + p2.y, p.z + p2.z}
}

func (p Point) Diff(p2 Point) Point {
	return Point{p.x - p2.x, p.y - p2.y, p.z - p2.z}
}

type Scanner struct {
	Points             []Point
	Pos                Point
	RelativeToScanner0 bool
}

type Axis int

const (
	X_AXIS Axis = iota
	Y_AXIS
	Z_AXIS
)

type AxisFlip int

const (
	NO_FLIP AxisFlip = 0
	X_FLIP  AxisFlip = 1 << iota
	Y_FLIP
	Z_FLIP

	XY_FLIP  = X_FLIP | Y_FLIP
	XZ_FLIP  = X_FLIP | Z_FLIP
	YZ_FLIP  = Y_FLIP | Z_FLIP
	XYZ_FLIP = XY_FLIP | Z_FLIP
)

type Orientation struct {
	Order [3]Axis
	Flip  AxisFlip
}

var Orientations = []Orientation{
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, NO_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, X_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, Y_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, Z_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, XY_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, XZ_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, YZ_FLIP},
	{[3]Axis{X_AXIS, Y_AXIS, Z_AXIS}, XYZ_FLIP},

	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, NO_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, X_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, Y_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, Z_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, XY_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, XZ_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, YZ_FLIP},
	{[3]Axis{X_AXIS, Z_AXIS, Y_AXIS}, XYZ_FLIP},

	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, NO_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, X_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, Y_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, Z_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, XY_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, XZ_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, YZ_FLIP},
	{[3]Axis{Y_AXIS, X_AXIS, Z_AXIS}, XYZ_FLIP},

	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, NO_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, X_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, Y_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, Z_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, XY_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, XZ_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, YZ_FLIP},
	{[3]Axis{Y_AXIS, Z_AXIS, X_AXIS}, XYZ_FLIP},

	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, NO_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, X_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, Y_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, Z_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, XY_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, XZ_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, YZ_FLIP},
	{[3]Axis{Z_AXIS, X_AXIS, Y_AXIS}, XYZ_FLIP},

	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, NO_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, X_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, Y_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, Z_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, XY_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, XZ_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, YZ_FLIP},
	{[3]Axis{Z_AXIS, Y_AXIS, X_AXIS}, XYZ_FLIP},
}

func getAxis(p Point, a Axis) int {
	switch a {
	case X_AXIS:
		return p.x
	case Y_AXIS:
		return p.y
	case Z_AXIS:
		return p.z
	}
	return 0
}

func orientScanner(points []Point, o Orientation) []Point {
	orientedPoints := make([]Point, len(points))
	copy(orientedPoints, points)

	if o == Orientations[0] {
		return orientedPoints
	}

	for i := range orientedPoints {
		point := orientedPoints[i]
		point.x = getAxis(orientedPoints[i], o.Order[0])
		point.y = getAxis(orientedPoints[i], o.Order[1])
		point.z = getAxis(orientedPoints[i], o.Order[2])

		if o.Flip&X_FLIP > 0 {
			point.x = -point.x
		}
		if o.Flip&Y_FLIP > 0 {
			point.y = -point.y
		}
		if o.Flip&Z_FLIP > 0 {
			point.z = -point.z
		}

		orientedPoints[i] = point
	}
	return orientedPoints
}

func countCopies(masterPointsList map[Point]struct{}, src []Point) int {
	count := 0
	for i := 0; i < len(src); i++ {
		for p := range masterPointsList {
			if src[i] == p {
				count++
			}
		}
	}
	return count
}

func canMap(masterPointsList map[Point]struct{}, from *Scanner, o Orientation) bool {
	orientedPoints := orientScanner(from.Points, o)
	for _, op := range orientedPoints {
		for toPoint := range masterPointsList {
			movedPoints := make([]Point, len(orientedPoints))
			copy(movedPoints, orientedPoints)

			diff := toPoint.Diff(op)
			for i := range movedPoints {
				movedPoints[i] = movedPoints[i].Add(diff)
			}

			if countCopies(masterPointsList, movedPoints) >= CommonPointsNeeded {
				from.Points = movedPoints
				from.Pos = diff
				from.RelativeToScanner0 = true
				return true
			}
		}
	}
	return false
}

func scannersSolved(scanners []Scanner) bool {
	for i := 0; i < len(scanners); i++ {
		if !scanners[i].RelativeToScanner0 {
			return false
		}
	}
	return true
}

func addToMasterList(masterList map[Point]struct{}, points []Point) {
	for _, p := range points {
		masterList[p] = struct{}{}
	}
}

func solve(scanners []Scanner) {
	masterPointsList := make(map[Point]struct{})
	addToMasterList(masterPointsList, scanners[0].Points)
	for !scannersSolved(scanners) {
	scannerLoop:
		for j := 0; j < len(scanners); j++ {
			if scanners[j].RelativeToScanner0 {
				continue
			}
			for _, orientation := range Orientations {
				if canMap(masterPointsList, &scanners[j], orientation) {
					fmt.Println(j, scanners[j].Pos)
					addToMasterList(masterPointsList, scanners[j].Points)
					continue scannerLoop
				}
			}
		}
	}
}

func countUniqueBeacons(scanners []Scanner) int {
	m := make(map[Point]struct{})

	for i := 0; i < len(scanners); i++ {
		for _, p := range scanners[i].Points {
			m[p] = struct{}{}
		}
	}
	return len(m)
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	scanners := []Scanner{}
	for _, scannerInfo := range strings.Split(string(input), "\n\n") {
		pointsInfo := strings.Split(scannerInfo, "\n")[1:]
		points := []Point{}
		for _, point := range pointsInfo {
			axis := strings.Split(point, ",")
			x, _ := strconv.Atoi(axis[0])
			y, _ := strconv.Atoi(axis[1])
			z, _ := strconv.Atoi(axis[2])
			points = append(points, Point{x, y, z})
		}
		scanners = append(scanners, Scanner{Points: points})
	}
	scanners[0].RelativeToScanner0 = true

	solve(scanners)
	fmt.Println(countUniqueBeacons(scanners))
}
