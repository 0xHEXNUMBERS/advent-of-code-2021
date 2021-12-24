package main

import "testing"

func TestMagnitude(t *testing.T) {
	testCases := []struct {
		Input  string
		Output int
	}{
		{
			Input:  "[[9,1],[1,9]]",
			Output: 129,
		},
		{
			Input:  "[[1,2],[[3,4],5]]",
			Output: 143,
		},
		{
			Input:  "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
			Output: 1384,
		},
		{
			Input:  "[[[[1,1],[2,2]],[3,3]],[4,4]]",
			Output: 445,
		},
		{
			Input:  "[[[[3,0],[5,3]],[4,4]],[5,5]]",
			Output: 791,
		},
		{
			Input:  "[[[[5,0],[7,4]],[5,5]],[6,6]]",
			Output: 1137,
		},
		{
			Input:  "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
			Output: 3488,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Magnitude test case: %s", testCase.Input)
		num := parseSnailfishNumber(testCase.Input)
		mag := num.Magnitude()
		if mag != testCase.Output {
			t.Errorf("Magnitude of %s failed -- got: %d | expected: %d", testCase.Input, mag, testCase.Output)
		}
	}
}

func TestMagnitudeTestCase(t *testing.T) {
	input := "[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]"
	output := 4140

	num := parseSnailfishNumber(input)
	mag := num.Magnitude()

	if output != mag {
		t.Errorf("TestCase Magnitude failed -- got: %d | expected: %d", mag, output)
	}
}
