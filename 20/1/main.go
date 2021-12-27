package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ROUNDS = 2

type Point struct {
	y int
	x int
}

func (p Point) Add(p2 Point) Point {
	return Point{p.y + p2.y, p.x + p2.x}
}

type Image struct {
	mapping       []bool
	image         map[Point]bool
	outsidePixels bool
}

func makeImage(mapStr, imageStr string) Image {
	var im Image
	im.mapping = make([]bool, len(mapStr))
	for i := range mapStr {
		im.mapping[i] = mapStr[i] == '#'
	}

	imageRows := strings.Split(imageStr, "\n")
	im.image = make(map[Point]bool, len(imageRows))
	for i, rows := range imageRows {
		for j := range rows {
			im.image[Point{i, j}] = imageRows[i][j] == '#'
		}
	}
	return im
}

func enhancePixel(im *Image, p Point) bool {
	diffs := []Point{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1},
		{0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	}

	index := 0
	for _, d := range diffs {
		p := p.Add(d)
		index <<= 1
		if isOn, ok := im.image[p]; (isOn && ok) || (!ok && im.outsidePixels) {
			index |= 1
		}
	}
	return im.mapping[index]
}

func enhance(im *Image) {
	newImage := make(map[Point]bool)

	maxX, maxY, minX, minY := 0, 0, 0, 0
	for p := range im.image {
		newImage[p] = enhancePixel(im, p)

		if maxY < p.y {
			maxY = p.y
		} else if p.y < minY {
			minY = p.y
		}
		if maxX < p.x {
			maxX = p.x
		} else if p.x < minX {
			minX = p.x
		}
	}

	for y := minY - 1; y <= maxY+1; y++ {
		pMin := Point{y, minX - 1}
		pMax := Point{y, maxX + 1}
		newImage[pMin] = enhancePixel(im, pMin)
		newImage[pMax] = enhancePixel(im, pMax)
	}
	for x := minX - 1; x <= maxX+1; x++ {
		pMin := Point{minY - 1, x}
		pMax := Point{maxY + 1, x}
		newImage[pMin] = enhancePixel(im, pMin)
		newImage[pMax] = enhancePixel(im, pMax)
	}

	if !im.outsidePixels && im.mapping[0] {
		im.outsidePixels = !im.outsidePixels
	} else if im.outsidePixels && !im.mapping[511] {
		im.outsidePixels = !im.outsidePixels
	}

	im.image = newImage
}

func countOnPixels(im Image) int {
	count := 0
	for _, isOn := range im.image {
		if isOn {
			count += 1
		}
	}
	return count
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	args := strings.Split(string(input), "\n\n")
	image := makeImage(args[0], args[1])

	for i := 0; i < ROUNDS; i++ {
		enhance(&image)
	}
	fmt.Println(countOnPixels(image))
}
