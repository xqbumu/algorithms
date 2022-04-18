package question

import "testing"

// https://leetcode-cn.com/problems/search-insert-position/
func searchInsert(nums []int, target int) int {
	if len(nums) == 0 || nums[0] > target {
		return 0
	}
	if nums[len(nums)-1] < target {
		return len(nums)
	}

	i, j := 0, len(nums)
	m := (i + j) / 2
	for {
		if nums[m] >= target {
			j = m
		} else if nums[m] <= target {
			i = m + 1
		}
		if m == (i+j)/2 {
			break
		}
		m = (i + j) / 2
	}

	res := m
	for res-1 >= 0 && nums[res-1] >= target {
		res--
	}
	for res+1 < len(nums) && nums[res+1] <= target {
		res++
	}

	return res
}

func Test_searchInsert(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{[]int{1, 3, 5, 6}, 5}, 2},
		{``, args{[]int{1, 3, 5, 6}, 2}, 1},
		{``, args{[]int{0, 1, 3, 3, 5, 6}, 2}, 2},
		{``, args{[]int{1, 3, 5, 6}, 7}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchInsert(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("searchInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}
