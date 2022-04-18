package main

import (
	"testing"
)

const lb, rb = '(', ')'

type statement struct {
	since    int
	until    int
	children []*statement
}

// https://leetcode-cn.com/problems/longest-valid-parentheses/
func longestValidParentheses(s string) int {
	var res, val, offset int
	for i := 0; i < len(s); {
		val, offset = singleLongest(s[i:])
		if val > res {
			res = val
			i += offset
			continue
		}
		i++
	}
	return res
}

func singleLongest(s string) (int, int) {
	var res, pos, val, offset int
	for pos < len(s) {
		if s[pos] == '(' && pos+1 < len(s) {
			if s[pos+1] == ')' {
				res += 2
				pos += 2
				continue
			}

			pos++
			val, offset = singleLongest(s[pos:])
			if pos+offset < len(s) && s[pos+offset] == ')' {
				res += val + 2
				pos += offset + 1
			} else if val > res {
				res = val
				break
			} else {
				break
			}
		} else {
			break
		}
	}
	return res, pos
}

func Test_longestValidParentheses(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{``, args{"(()"}, 2},
		{``, args{"(())"}, 4},
		{``, args{")()())"}, 4},
		{``, args{"()(("}, 2},
		{``, args{""}, 0},
		{``, args{"()(()"}, 2},
		{``, args{"()(()("}, 2},
		{``, args{"()(()()("}, 4},
		{``, args{"()()(()("}, 4},
		{``, args{"(())(()("}, 4},
		{``, args{"()((())("}, 4},
		{`lc18`, args{")()())()()("}, 4},
		{`lc18`, args{")()()))()()("}, 4},
		{`lc18`, args{")()()))(())()("}, 6},
		{`lc19`, args{"))(()(("}, 2},
		{`lc19`, args{"))))((()(("}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestValidParentheses(tt.args.s); got != tt.want {
				t.Errorf("longestValidParentheses() = %v, want %v", got, tt.want)
			}
		})
	}
}
