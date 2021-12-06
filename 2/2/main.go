package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Submarine struct {
	y   int
	x   int
	aim int
}

func (s *Submarine) Execute(command string, i int) {
	switch command {
	case "forward":
		s.x += i
		s.y += s.aim * i
	case "down":
		s.aim += i
	case "up":
		s.aim -= i
	}
}

func main() {
	sub := Submarine{}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(input), "\n") {
		if line == "" {
			continue
		}

		args := strings.Split(line, " ")
		i, _ := strconv.Atoi(args[1])
		sub.Execute(args[0], i)
	}
	fmt.Println(sub.y * sub.x)
}
