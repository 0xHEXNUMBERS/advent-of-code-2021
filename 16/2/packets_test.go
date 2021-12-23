package main

import "testing"

func TestPacketValue(t *testing.T) {
	testCases := []struct {
		Input  string
		Output uint64
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}

	for _, testCase := range testCases {
		t.Logf("Running test case %s", testCase.Input)
		out := packetValue(testCase.Input)
		if out != testCase.Output {
			t.Errorf(
				"Test Case %s fails: got: %d | expected: %d",
				testCase.Input, out, testCase.Output)
		}
	}
}
