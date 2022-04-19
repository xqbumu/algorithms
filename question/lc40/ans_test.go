package question

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://leetcode-cn.com/problems/combination-sum-ii/

func combinationSum2(candidates []int, target int) [][]int {
	candidates = findValid(candidates, target)
	sort.Ints(candidates)

	res, subRes := [][]int{}, [][]int{}
	for i := 0; i < len(candidates); i++ {
		if target-candidates[i] == 0 {
			var found bool
			for j := 0; j < len(res); j++ {
				if exist(res, []int{candidates[i]}) {
					found = true
					break
				}
			}
			if !found {
				res = append(res, []int{candidates[i]})
			}
		} else {
			subRes = combinationSum2(candidates[i+1:], target-candidates[i])
			if len(subRes) > 0 {
				for j := 0; j < len(subRes); j++ {
					subRes[j] = append(subRes[j], candidates[i])
					if exist(res, subRes[j]) {
						continue
					}
					res = append(res, subRes[j])
				}
			}
		}
	}

	return res
}

func exist(a [][]int, b []int) bool {
	for i := 0; i < len(a); i++ {
		if compare(a[i], b) {
			return true
		}
	}
	return false
}

func compare(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
	sum := 0
	sums := map[int]int{}
	nums := make([]int, 0, len(candidates))
	for i := 0; i < len(candidates); i++ {
		if candidates[i] <= target && sums[candidates[i]] < target {
			sums[candidates[i]] += candidates[i]
			sum += candidates[i]
			nums = append(nums, candidates[i])
		}
	}

	if sum < target {
		return []int{}
	}

	return nums
}

func Test_combinationSum2(t *testing.T) {
	type args struct {
		candidates []int
		target     int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{``, args{[]int{10, 1, 2, 7, 6, 1, 5}, 8}, [][]int{
			{1, 1, 6},
			{1, 2, 5},
			{1, 7},
			{2, 6},
		}},
		{``, args{[]int{2, 5, 2, 1, 2}, 5}, [][]int{
			{1, 2, 2},
			{5},
		}},
		{`t131`, args{[]int{1, 1}, 1}, [][]int{
			{1},
		}},
		{`t131`, args{[]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 30}, [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combinationSum2(tt.args.candidates, tt.args.target)
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
