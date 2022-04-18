package question

import "testing"

// https://leetcode-cn.com/problems/search-in-rotated-sorted-array/
func search(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	for i := 0; i < len(nums); i++ {
		if nums[i] == target {
			return i
		}
	}
	return -1
}

func Test_search(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{[]int{4, 5, 6, 7, 0, 1, 2}, 0}, 4},
		{``, args{[]int{4, 5, 6, 7, 0, 1, 2}, 3}, -1},
		{``, args{[]int{1}, 0}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}
