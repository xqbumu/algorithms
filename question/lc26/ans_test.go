package main

import (
	"reflect"
	"testing"
)

// https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array/
func removeDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}

	count := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			continue
		}
		count++
		nums[count] = nums[i]
	}
	return count + 1
}

func Test_removeDuplicates(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
		nums []int
	}{
		{``, args{[]int{1, 1, 2}}, 2, []int{1, 2}},
		{``, args{[]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}}, 5, []int{0, 1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.args.nums)

			if got != tt.want {
				t.Errorf("removeDuplicates() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.args.nums[:got], tt.nums[:got]) {
				t.Errorf("removeDuplicates() = %v, want %v", tt.args.nums[:got], tt.nums[:got])
			}
		})
	}
}
