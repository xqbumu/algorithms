package question

import (
	"testing"
)

// https://leetcode-cn.com/problems/first-missing-positive/
func firstMissingPositive(nums []int) int {
	cnt := len(nums)
	if cnt == 0 {
		return 1
	}

	for i := 1; i <= cnt; i++ {
		p := i - 1
		if nums[p] == i {
			continue
		}
		if nums[p] > cnt || nums[p] < 1 {
			nums[p] = -1
			continue
		}
		swap(nums, i)
	}

	for i := 1; i <= cnt; i++ {
		if nums[i-1] != i {
			return i
		}
	}

	return nums[cnt-1] + 1
}

func swap(nums []int, i int) {
	if i-1 < 0 {
		return
	}
	a := nums[i-1]
	if i == a {
		return
	}
	b := nums[a-1]
	if a == b {
		b = -1
	}
	nums[a-1] = a
	nums[i-1] = b
	if b > 0 {
		swap(nums, i)
	}
}

func Test_firstMissingPositive(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{[]int{2147483647, 1}}, 2},
		{``, args{[]int{2147483647, 2, 5}}, 1},
		{``, args{[]int{2147483647}}, 1},
		{``, args{[]int{0}}, 1},
		{``, args{[]int{1}}, 2},
		{``, args{[]int{1, 1}}, 2},
		{``, args{[]int{1, 2}}, 3},
		{``, args{[]int{2, 1}}, 3},
		{``, args{[]int{1, 2, 0}}, 3},
		{``, args{[]int{3, 4, 2, -1}}, 1},
		{``, args{[]int{3, 4, -1, 1}}, 2},
		{``, args{[]int{7, 8, 9, 11, 12}}, 1},
		{``, args{[]int{3, 4, 2, 2, 9, 5, 1, 13, -1, -1, -4, 3, 15, -10, 6, 10}}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstMissingPositive(tt.args.nums); got != tt.want {
				t.Errorf("firstMissingPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}
