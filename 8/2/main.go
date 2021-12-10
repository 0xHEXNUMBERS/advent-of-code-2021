package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type DigitDeducer struct {
	OnesPlaceholder string
	OnesSegments    [2]byte

	FoursPlaceholder string
	FoursSegments    [2]byte

	SevensPlaceholder string
	SevensSegment     byte
}

func (d *DigitDeducer) DeductBaseDigits() {
	//Grab ones segments
	ones := d.OnesPlaceholder
	d.OnesSegments[0] = ones[0]
	d.OnesSegments[1] = ones[1]

	//Grab fours segments
	fours := d.FoursPlaceholder
	foursIndex := 0
	for i := 0; i < len(fours); i++ {
		if fours[i] != ones[0] && fours[i] != ones[1] {
			d.FoursSegments[foursIndex] = fours[i]
			foursIndex++
		}
	}

	//Grab seven segments
	sevens := d.SevensPlaceholder
	for i := 0; i < len(sevens); i++ {
		if sevens[i] != ones[0] &&
			sevens[i] != ones[1] {
			d.SevensSegment = sevens[i]
			break
		}
	}
}

func (d *DigitDeducer) Deduce5Segments(segments string) int {
	//5 uses both 4 segments; 2 and 3 do not
	//3 uses both 1 segments; 2 does not
	count4Segments := 0
	count1Segments := 0
	for i := 0; i < len(segments); i++ {
		if segments[i] == d.FoursSegments[0] ||
			segments[i] == d.FoursSegments[1] {
			count4Segments++
		}

		if segments[i] == d.OnesSegments[0] ||
			segments[i] == d.OnesSegments[1] {
			count1Segments++
		}
	}
	if count4Segments == 2 {
		return 5
	} else if count1Segments == 2 {
		return 3
	} else {
		return 2
	}
}

func (d *DigitDeducer) Deduce6Segments(segments string) int {
	//6 uses one 1 segments; 9 and 0 do not
	//9 uses both 4 segments; 0 does not
	count4Segments := 0
	count1Segments := 0
	for i := 0; i < len(segments); i++ {
		if segments[i] == d.OnesSegments[0] ||
			segments[i] == d.OnesSegments[1] {
			count1Segments++
		}

		if segments[i] == d.FoursSegments[0] ||
			segments[i] == d.FoursSegments[1] {
			count4Segments++
		}
	}
	if count1Segments == 1 {
		return 6
	} else if count4Segments == 2 {
		return 9
	} else {
		return 0
	}
}

func (d *DigitDeducer) DeduceDigit(segments string) int {
	switch len(segments) {
	case 2:
		return 1
	case 3:
		return 7
	case 4:
		return 4
	case 5:
		return d.Deduce5Segments(segments)
	case 6:
		return d.Deduce6Segments(segments)
	case 7:
		return 8
	default:
		return -100000000 //Easy debugging
	}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last \n

	lines := strings.Split(string(input), "\n")
	sum := 0
	for _, line := range lines {
		args := strings.Split(line, " | ")
		inputSegments := strings.Split(args[0], " ")
		outputSegments := strings.Split(args[1], " ")

		d := DigitDeducer{}
		for _, segments := range inputSegments {
			switch len(segments) {
			case 2:
				d.OnesPlaceholder = segments
			case 3:
				d.SevensPlaceholder = segments
			case 4:
				d.FoursPlaceholder = segments
			}
		}
		d.DeductBaseDigits()

		outputNumber := 0
		for _, segments := range outputSegments {
			outputNumber *= 10
			outputNumber += d.DeduceDigit(segments)
		}
		sum += outputNumber
	}
	fmt.Println(sum)
}
