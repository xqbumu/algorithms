package main

import (
	"reflect"
	"testing"
)

// [TODO] https://leetcode-cn.com/problems/next-permutation/
func nextPermutation(nums []int) {

}

func Test_nextPermutation(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{``, args{[]int{1, 2, 3}}, []int{1, 3, 2}},
		{``, args{[]int{2, 3, 1}}, []int{3, 1, 2}},
		{``, args{[]int{3, 2, 1}}, []int{1, 2, 3}},
		{``, args{[]int{1, 1, 5}}, []int{1, 5, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextPermutation(tt.args.nums)
			if !reflect.DeepEqual(tt.args.nums, tt.want) {
				t.Errorf("removeElement() = %v, want %v", tt.args.nums, tt.want)
			}
		})
	}
}
