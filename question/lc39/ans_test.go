package question

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://leetcode-cn.com/problems/combination-sum/
func combinationSum(candidates []int, target int) [][]int {
	candidates = findValid(candidates, target)

	res := [][]int{}
	for i := 0; i < len(candidates); i++ {
		if target-candidates[i] == 0 {
			res = append(res, []int{candidates[i]})
		} else {
			subRes := combinationSum(candidates[i:], target-candidates[i])
			if len(subRes) > 0 {
				subRes = crossResult([][]int{{candidates[i]}}, subRes)
				for i := 0; i < len(subRes); i++ {
					res = append(res, subRes[i])
				}
			}
		}
	}

	return res
}

func crossResult(a, b [][]int) [][]int {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	res := make([][]int, 0, len(a)*len(b))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			res = append(res, append(append([]int{}, a[i]...), b[j]...))
		}
	}
	return res
}

func findValid(candidates []int, target int) []int {
	nums := make([]int, 0, len(candidates))
	for i := 0; i < len(candidates); i++ {
		if candidates[i] <= target {
			nums = append(nums, candidates[i])
		}
	}
	return nums
}

func Test_combinationSum(t *testing.T) {
	type args struct {
		candidates []int
		target     int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{``, args{[]int{1}, 1}, [][]int{{1}}},
		{``, args{[]int{1, 2}, 2}, [][]int{{1, 1}, {2}}},
		{``, args{[]int{2, 3, 6, 7}, 7}, [][]int{{2, 2, 3}, {7}}},
		{``, args{[]int{2, 3, 5}, 8}, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}},
		{``, args{[]int{2}, 1}, [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combinationSum(tt.args.candidates, tt.args.target)
			sortResult(tt.want)
			sortResult(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func sortResult(data [][]int) {
	for i := 0; i < len(data); i++ {
		sort.Ints(data[i])
	}
}
