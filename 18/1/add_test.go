package main

import "testing"

func TestAdd(t *testing.T) {
	testCases := []struct {
		Input  [2]string
		Output string
	}{
		{
			Input: [2]string{
				"[[[[4,3],4],4],[7,[[8,4],9]]]",
				"[1,1]",
			},
			Output: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
	}

	for _, testCase := range testCases {
		t.Logf("Reduce Test Case: %s", testCase.Input)
		sfn1 := parseSnailfishNumber(testCase.Input[0])
		sfn2 := parseSnailfishNumber(testCase.Input[1])
		sfnSum := sfn1.Add(sfn2)
		if sfnSum.String() != testCase.Output {
			t.Errorf("Reduce failed on %s + %s -- got: %s | expected: %s", testCase.Input[0], testCase.Input[1], sfnSum, testCase.Output)
		}
	}
}

func TestAddTestCase(t *testing.T) {
	input := []string{
		"[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]",
		"[[[5,[2,8]],4],[5,[[9,9],0]]]",
		"[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]",
		"[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]",
		"[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]",
		"[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]",
		"[[[[5,4],[7,7]],8],[[8,3],8]]",
		"[[9,3],[[9,9],[6,[4,9]]]]",
		"[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]",
		"[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]",
	}
	output := "[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]"

	sum := parseSnailfishNumber(input[0])
	for i := 1; i < len(input); i++ {
		num := parseSnailfishNumber(input[i])
		sum = sum.Add(num)
	}
	if output != sum.String() {
		t.Errorf("TestCase Add failed -- got: %s | expected: %s", sum, output)
	}
}
