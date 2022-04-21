package question

import (
	"testing"
)

// https://leetcode-cn.com/problems/trapping-rain-water/
func trap(height []int) int {
	if len(height) < 3 {
		return 0
	}

	var res, lm, rm, s int
	for i := 1; i < len(height)-1; i++ {
		lm, rm = 0, 0
		for j := 0; j < len(height); j++ {
			if j < i && height[j] > lm {
				lm = height[j]
			}
			if j > i && height[j] > rm {
				rm = height[j]
			}
		}

		s = lm
		if rm < s {
			s = rm
		}

		if s < height[i] {
			continue
		}

		res += (s - height[i])
	}

	return res
}

// 0
// 0 0             0
// 0 0 0     0     0
// 0 0 0 0   0     0
// 0 0 0 0 0 0     0     0
// 0 0 0 0 0 0 0   0   0 0
// 0 0 0 0 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0 0 0 0 0
// 0 0 0 0 0 0 0 0 0 0 0 0
// 90807060507040308030405

//        0
//    0   00 0  4
//  0 00 000000 2
// 010210132121
// 0123456789

//      0
// 0    0
// 0  0 0
// 00 000
// 00 000
// 420325

func Test_trap(t *testing.T) {
	type args struct {
		height []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{[]int{0, 2, 0}}, 0},
		{``, args{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}}, 6},
		{``, args{[]int{4, 2, 0, 3, 2, 5}}, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trap(tt.args.height); got != tt.want {
				t.Errorf("trap() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func trapSlow(height []int) int {
// 	if len(height) < 3 {
// 		return 0
// 	}

// 	min, max := height[0], height[0]
// 	for i := 1; i < len(height); i++ {
// 		if height[i] < min {
// 			min = height[i]
// 		}
// 		if height[i] > max {
// 			max = height[i]
// 		}
// 	}

// 	var res, prev int
// 	for i := min; i < max; i++ {
// 		prev = -1
// 		for j := 0; j < len(height); j++ {
// 			if height[j] <= i {
// 				continue
// 			}
// 			if prev == -1 {
// 				prev = j
// 				continue
// 			}
// 			res += j - prev - 1
// 			prev = j
// 		}
// 	}

// 	return res
// }
