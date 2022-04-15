package main

import (
	"reflect"
	"sort"
	"testing"
)

// https://leetcode-cn.com/problems/remove-element/
func removeElement(nums []int, val int) int {
	count := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] == val {
			continue
		}
		nums[count] = nums[i]
		count++
	}
	return count
}

func Test_removeElement(t *testing.T) {
	type args struct {
		nums []int
		val  int
	}
	tests := []struct {
		name string
		args args
		want int
		nums []int
	}{
		{``, args{[]int{3, 2, 2, 3}, 3}, 2, []int{2, 2}},
		{``, args{[]int{0, 1, 2, 2, 3, 0, 4, 2}, 2}, 5, []int{0, 1, 4, 0, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeElement(tt.args.nums, tt.args.val)

			if got != tt.want {
				t.Errorf("removeElement() = %v, want %v", got, tt.want)
			}

			left := sort.IntSlice(tt.args.nums[:got])
			left.Sort()
			right := sort.IntSlice(tt.nums[:got])
			right.Sort()
			if !reflect.DeepEqual(left, right) {
				t.Errorf("removeElement() = %v, want %v", left, right)
			}
		})
	}
}
