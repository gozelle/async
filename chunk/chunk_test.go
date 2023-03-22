package chunk

import (
	"github.com/gozelle/testify/require"
	"testing"
)

func TestInt64Ranges(t *testing.T) {
	
	type testCase struct {
		Nums   []int64
		Result [][]int64
		Error  bool
	}
	
	testCases := []testCase{
		{Nums: []int64{0, 0, 3}, Result: [][]int64{{0, 0}}},
		{Nums: []int64{0, 1, 3}, Result: [][]int64{{0, 1}}},
		{Nums: []int64{0, 9, 3}, Result: [][]int64{{0, 2}, {3, 5}, {6, 8}, {9, 9}}},
		{Nums: []int64{1, 5, 2}, Result: [][]int64{{1, 2}, {3, 4}, {5, 5}}},
		{Nums: []int64{1, 8, 2}, Result: [][]int64{{1, 2}, {3, 4}, {5, 6}, {7, 8}}},
		{Nums: []int64{0, 3, 5}, Result: [][]int64{{0, 3}}},
		{Nums: []int64{0, 101, 50}, Result: [][]int64{{0, 49}, {50, 99}, {100, 101}}},
		{Nums: []int64{3, 9, 4}, Result: [][]int64{{3, 6}, {7, 9}}},
		{Nums: []int64{-10, 10, 5}, Result: [][]int64{{-10, -6}, {-5, -1}, {0, 4}, {5, 9}, {10, 10}}},
		{Nums: []int64{-10, 10, -1}, Error: true},
		{Nums: []int64{-10, 10, 0}, Error: true},
	}
	
	for _, v := range testCases {
		r, err := SplitInt64s(v.Nums[0], v.Nums[1], v.Nums[2])
		if v.Error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		t.Log(v.Nums, r)
		require.Equal(t, len(r), len(v.Result))
		for i, vv := range r {
			require.Equal(t, v.Result[i][0], vv.Begin)
			require.Equal(t, v.Result[i][1], vv.End)
		}
	}
}
