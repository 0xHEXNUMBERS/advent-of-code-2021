package main

import "testing"

func TestParse(t *testing.T) {
	testCases := []struct {
		Input  string
		Output *SnailfishNumber
	}{
		{
			Input: "[[9,1],[1,9]]",
			Output: &SnailfishNumber{
				Left: &SnailfishNumber{
					Left:  &SnailfishNumber{Number: 9},
					Right: &SnailfishNumber{Number: 1},
				},
				Right: &SnailfishNumber{
					Left:  &SnailfishNumber{Number: 1},
					Right: &SnailfishNumber{Number: 9},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Parse Test Case: %s", testCase.Input)
		out := parseSnailfishNumber(testCase.Input)
		if out.String() != testCase.Output.String() {
			t.Errorf("Parsing %s failed -- got: %s | expected: %s", testCase.Input, *out, *testCase.Output)
		}
	}
}
