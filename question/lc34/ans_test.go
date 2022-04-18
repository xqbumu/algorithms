package question

import (
	"reflect"
	"testing"
)

// https://leetcode-cn.com/problems/find-first-and-last-position-of-element-in-sorted-array/
func searchRange(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{-1, -1}
	}

	c, i, j := 0, 0, len(nums)
	m := (i + j) / 2
	for nums[m] != target {
		c++
		if nums[m] > target {
			j = m
		} else {
			i = m
		}
		if m == (i+j)/2 {
			break
		}
		m = (i + j) / 2
	}
	if nums[m] != target {
		return []int{-1, -1}
	}

	i, j = m, m
	for i-1 >= 0 && nums[i-1] == target {
		i--
	}
	for j+1 < len(nums) && nums[j+1] == target {
		j++
	}

	return []int{i, j}
}

func Test_searchRange(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{``, args{[]int{-1, 1, 2, 2, 4, 5, 6, 7, 8, 9, 10}, 2}, []int{2, 3}},
		{``, args{[]int{5, 7, 7, 8, 8, 10}, 8}, []int{3, 4}},
		{``, args{[]int{5, 7, 7, 8, 8, 10}, 6}, []int{-1, -1}},
		{``, args{[]int{}, 0}, []int{-1, -1}},
		{``, args{[]int{1}, 1}, []int{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchRange(tt.args.nums, tt.args.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
