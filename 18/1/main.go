package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type SnailfishNumber struct {
	Left   *SnailfishNumber
	Right  *SnailfishNumber
	Number int
}

func (s SnailfishNumber) IsPair() bool {
	return s.Left != nil && s.Right != nil
}

func (s SnailfishNumber) String() string {
	if s.IsPair() {
		return "[" + s.Left.String() + "," + s.Right.String() + "]"
	}
	return fmt.Sprintf("%d", s.Number)
}

func (s SnailfishNumber) Magnitude() int {
	if s.IsPair() {
		return 3*s.Left.Magnitude() + 2*s.Right.Magnitude()
	}
	return s.Number
}

func (s *SnailfishNumber) ExplosionAddRightValue(value int) {
	if s.Left != nil {
		s.Left.ExplosionAddRightValue(value)
	} else {
		s.Number += value
	}
}

func (s *SnailfishNumber) ExplosionAddLeftValue(value int) {
	if s.Right != nil {
		s.Right.ExplosionAddLeftValue(value)
	} else {
		s.Number += value
	}
}

func (s *SnailfishNumber) SearchForExplosion(depth int) (int, int, bool, bool, bool) {
	if depth == 4 && s.IsPair() {
		return s.Left.Number, s.Right.Number, true, true, true
	}

	if s.Left != nil {
		l, r, lok, rok, ok := s.Left.SearchForExplosion(depth + 1)
		if lok && rok {
			s.Left = &SnailfishNumber{Number: 0}
		}
		if rok {
			s.Right.ExplosionAddRightValue(r)
			rok = false
		}
		if ok {
			return l, r, lok, rok, ok
		}
	}
	if s.Right != nil {
		l, r, lok, rok, ok := s.Right.SearchForExplosion(depth + 1)
		if lok && rok {
			s.Right = &SnailfishNumber{Number: 0}
		}
		if lok {
			s.Left.ExplosionAddLeftValue(l)
			lok = false
		}
		if ok {
			return l, r, lok, rok, ok
		}
	}
	return 0, 0, false, false, false
}

func (s *SnailfishNumber) SearchForSplit() bool {
	if s.Left.IsPair() {
		found := s.Left.SearchForSplit()
		if found {
			return true
		}
	} else {
		num := s.Left.Number
		if num >= 10 {
			s.Left = &SnailfishNumber{
				Left: &SnailfishNumber{
					Number: num / 2,
				},
				Right: &SnailfishNumber{
					Number: (num / 2) + (num % 2),
				},
			}
			return true
		}
	}
	if s.Right.IsPair() {
		found := s.Right.SearchForSplit()
		if found {
			return true
		}
	} else {
		num := s.Right.Number
		if num >= 10 {
			s.Right = &SnailfishNumber{
				Left: &SnailfishNumber{
					Number: num / 2,
				},
				Right: &SnailfishNumber{
					Number: (num / 2) + (num % 2),
				},
			}
			return true
		}
	}
	return false
}

func (s *SnailfishNumber) Reduce() {
	for {
		_, _, _, _, ok := s.SearchForExplosion(0)
		if ok {
			continue
		}
		ok = s.SearchForSplit()
		if !ok {
			return
		}
	}
}

func (s *SnailfishNumber) Add(s2 *SnailfishNumber) *SnailfishNumber {
	sSum := &SnailfishNumber{Left: s, Right: s2}
	sSum.Reduce()
	return sSum
}

func extractParts(str string) (left, right string) {
	count := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '[' {
			count++
		} else if str[i] == ']' {
			count--
		}
		if count == 0 {
			return str[:i+1], str[i+2:] //remove ,
		}
	}
	panic("unreachable")
}

func parseSnailfishNumber(str string) *SnailfishNumber {
	if str[0] != '[' {
		num, _ := strconv.Atoi(str)
		return &SnailfishNumber{Number: num}
	}

	str = str[1 : len(str)-1] //Remove starting '[' & ending ']'
	s := &SnailfishNumber{}

	left, right := extractParts(str)

	s.Left = parseSnailfishNumber(left)
	s.Right = parseSnailfishNumber(right)

	return s
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	numbers := strings.Split(string(input), "\n")
	sum := parseSnailfishNumber(numbers[0])
	for i := 1; i < len(numbers); i++ {
		num := parseSnailfishNumber(numbers[i])
		sum = sum.Add(num)
	}
	fmt.Println(sum.Magnitude())
}
