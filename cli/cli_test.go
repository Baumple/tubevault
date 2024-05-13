package cli_test

import (
	"testing"

	"github.com/baumple/watchvault/cli"
	"github.com/baumple/watchvault/utility"
)

type WindowingTest struct {
	cursor int
	width  int

	expectedStart int
	expectedEnd   int
}

var windowingTests = []WindowingTest{
	{0, 10, 0, 10},
	{1, 10, 0, 10},
	{2, 10, 0, 10},
	{3, 10, 0, 10},
	{4, 10, 0, 10},
	{5, 10, 0, 10},
	{6, 10, 1, 11},
	{7, 10, 2, 12},
}

func TestWindowing(t *testing.T) {
	for _, test := range windowingTests {
		res := cli.GetWindow(test.cursor, test.width)
		if res.Start != test.expectedStart {
			t.Fatalf("Wanted start %d, got start %d", test.expectedStart, res.Start)
		}
		if res.End != test.expectedEnd {
			t.Fatalf("Wanted end %d, got end %d", test.expectedEnd, res.End)
		}
	}
}

type TestSplit struct {
	inputS string
	inputI int

	expected []string
}

var splitTests = []TestSplit{
	{"aabbcc", 2, []string{"aa", "bb", "cc"}},
	{"aaabbbcc", 3, []string{"aaa", "bbb", "cc"}},
}

func TestSplitEveryN(t *testing.T) {
	for _, test := range splitTests {
		res := utility.SplitEveryN(test.inputS, test.inputI)

		if len(res) != len(test.expected) {
			t.Fatalf("Wanted %v, got %v", test.expected, res)
		}

		for i, out := range res {
			if out != test.expected[i] {
				t.Fatalf("Wanted %v, got %v", test.expected, res)
			}
		}
	}
}

