package tests

import (
	"testing"
)

// quicksort sorts a given video slice
func quicksort(lowerBounds int, upperBounds int, videos []int) {
	if lowerBounds >= upperBounds || lowerBounds < 0 {
		return
	}

	p := partition(lowerBounds, upperBounds, videos)
	quicksort(lowerBounds, p-1, videos)
	quicksort(p+1, upperBounds, videos)
}

// 1 3 2 4 -> next iteration (left, right of pivot has n = 1)
// 1 3 2 pivot = 2
//       pivpos = 1 since 3 (with piv pos of 1) is swapped with 2 THE PIVOT IS NOW AT PIVPOS
//
// 1 3 2 6 4 2 pivot 2 pivotpos + 1
// 1 3 2 6 4 2 pivot 2 pivotpos = 1 -> stays the same and swap since 2 < 3
// 1 2 3 6 4 2 pivot 2 pivotpos = 2 -> swap i, pivpos (3 <-> 2)
// 1 2 3 6 4 2 pivot 2 pivotpos = 2
// 1 2 3 6 4 2 pivot 2 pivotpos = 2 -> swap i, pivpos (3 <-> 2)
// 1 2 2 3 6 4
// two "lists" ->  [1 2] (2) [3 6 4]
// 1 2 pivot: 2 pivotpos = 0
// 3 6 4 pivot: 4 pos: 0
// 3 6 4 pivot: 4 pos: 1
// 3 6 4 pivot: 4 pos: 1
// 3 4 6 // this sappens after the loop exits

func partition(lowerBounds int, upperBounds int, videos []int) int {
	pivotElement := videos[upperBounds]
	pivotPos := lowerBounds
	for i := lowerBounds; i < upperBounds; i++ {
		if videos[i] <= pivotElement {
			swap(i, pivotPos, videos)
			pivotPos++
		}
	}

	swap(pivotPos, upperBounds, videos)
	return pivotPos
}

func swap[T any](a int, b int, nums []T) {
	temp := nums[a]
	nums[a] = nums[b]
	nums[b] = temp
}

type QuickSortTest struct {
	input    []int
	expected []int
}

var quicksortTests = []QuickSortTest{
	{[]int{1, 392, 2, 321, -5, 22}, []int{-5, 1, 2, 22, 321, 392}},
	{[]int{999, 6969, 420, 42069, 69420, 42}, []int{42, 420, 999, 6969, 42069, 69420}},
	{[]int{666, -666, -999, -6969, -420, -42069, -69420, -42}, []int{-69420, -42069, -6969, -999, -666, -420, -42, 666}},
}

func TestQuickSort(t *testing.T) {
        for _, test := range quicksortTests {
                input := test.input
                quicksort(0, len(input) - 1, input)
                for idx, exp := range test.expected {
                        if exp != input[idx] {
                                t.Fatalf("Got %v, wanted %v", input, test.expected)
                        }
                }
       }
}
