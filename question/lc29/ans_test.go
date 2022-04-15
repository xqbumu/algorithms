package main

import (
	"log"
	"testing"
)

// https://leetcode-cn.com/problems/divide-two-integers/
func divide(dividend int, divisor int) int {
	valCheck := func(val int) int {
		if val < -1<<31 || val > (1<<31-1) {
			return 1<<31 - 1
		}
		return val
	}

	if divisor == 0 {
		panic("divisor can not be zero")
	}

	var res = 0
	var positive = true

	if dividend < 0 {
		dividend = -dividend
		positive = !positive
	}
	if divisor < 0 {
		divisor = -divisor
		positive = !positive
	}

	switch divisor {
	case 0:
		panic("divisor can not be zero")
	case 1:
		res = dividend
	case 2:
		res = dividend >> 1
	default:
		for dividend >= divisor {
			res++
			dividend -= divisor
		}
	}

	if !positive {
		return -res
	}

	return valCheck(res)
}

func TestD(t *testing.T) {
	a, b := 13, 3 // 1101 11 100
	a, b = 33, 3  // 100001, 11, 1011

	log.Println(a ^ b)
}

func Test_divide(t *testing.T) {
	type args struct {
		dividend int
		divisor  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{9, 3}, 3},
		{``, args{10, 3}, 3},
		{``, args{7, -3}, -2},
		{``, args{-2147483648, 2}, -1073741824},
		{``, args{-2147483648, -1}, 2147483647},
		{``, args{-2147483648, 1}, -2147483648},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := divide(tt.args.dividend, tt.args.divisor); got != tt.want {
				t.Errorf("divide() = %v, want %v", got, tt.want)
			}
		})
	}
}
