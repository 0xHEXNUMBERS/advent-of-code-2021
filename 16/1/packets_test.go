package main

import "testing"

func TestSumOfVersions(t *testing.T) {
	testCases := []struct {
		Input  string
		Output int
	}{
		{"D2FE28", 6},
		{"38006F45291200", 9},
		{"EE00D40C823060", 14},
		{"8A004A801A8002F478", 16},
		{"620080001611562C8802118E34", 12},
		{"C0015000016115A2E0802F182340", 23},
		{"A0016C880162017C3686B18A3D4780", 31},
	}

	for _, testCase := range testCases {
		t.Logf("Running test case %s", testCase.Input)
		out := versionSum(testCase.Input)
		if out != testCase.Output {
			t.Errorf(
				"Test Case %s fails: got: %d | expected: %d",
				testCase.Input, out, testCase.Output)
		}
	}
}
