package main

import "testing"

func TestExplode(t *testing.T) {
	testCases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "[[[[[9,8],1],2],3],4]",
			Output: "[[[[0,9],2],3],4]",
		},
		{
			Input:  "[7,[6,[5,[4,[3,2]]]]]",
			Output: "[7,[6,[5,[7,0]]]]",
		},
		{
			Input:  "[[6,[5,[4,[3,2]]]],1]",
			Output: "[[6,[5,[7,0]]],3]",
		},
		{
			Input:  "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			Output: "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
		},
		{
			Input:  "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
			Output: "[[3,[2,[8,0]]],[9,[5,[7,0]]]]",
		},
	}

	for _, testCase := range testCases {
		t.Logf("Explode Test Case: %s", testCase.Input)
		sfn := parseSnailfishNumber(testCase.Input)
		sfn.SearchForExplosion(0)
		if sfn.String() != testCase.Output {
			t.Errorf("Explosion failed on %s -- got: %s | expected: %s", testCase.Input, sfn, testCase.Output)
		}
	}
}
