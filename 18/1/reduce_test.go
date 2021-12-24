package main

import (
	"testing"
)

func TestReduce(t *testing.T) {
	testCases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]",
			Output: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
	}

	for _, testCase := range testCases {
		t.Logf("Reduce Test Case: %s", testCase.Input)
		sfn := parseSnailfishNumber(testCase.Input)
		sfn.Reduce()
		if sfn.String() != testCase.Output {
			t.Errorf("Reduce failed on %s -- got: %s | expected: %s", testCase.Input, sfn, testCase.Output)
		}
	}
}
