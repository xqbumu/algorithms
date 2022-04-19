package question

import (
	"fmt"
	"testing"
)

// https://leetcode-cn.com/problems/count-and-say/
// 1.     1
// 2.     11
// 3.     21
// 4.     1211
// 5.     111221
// 第一项是数字 1
// 描述前一项，这个数是 1 即 “ 一 个 1 ”，记作 "11"
// 描述前一项，这个数是 11 即 “ 二 个 1 ” ，记作 "21"
// 描述前一项，这个数是 21 即 “ 一 个 2 + 一 个 1 ” ，记作 "1211"
// 描述前一项，这个数是 1211 即 “ 一 个 1 + 一 个 2 + 二 个 1 ” ，记作 "111221"
func countAndSay(n int) string {
	return exec("1", n-1)
}

func exec(s string, n int) string {
	if n == 0 {
		return s
	}
	var (
		res  string
		cnt  int
		prev byte
	)
	for i := 0; i < len(s); i++ {
		if i == 0 {
			prev = s[i]
			cnt++
		} else if s[i] == prev {
			cnt++
		} else if s[i] != prev {
			res += fmt.Sprintf("%d%c", cnt, prev)
			prev = s[i]
			cnt = 1
		}
	}
	res += fmt.Sprintf("%d%c", cnt, prev)
	return exec(res, n-1)
}

func Test_countAndSay(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// {``, args{1}, "1"},
		// {``, args{2}, "11"},
		// {``, args{3}, "21"},
		{``, args{4}, "1211"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countAndSay(tt.args.n); got != tt.want {
				t.Errorf("countAndSay() = %v, want %v", got, tt.want)
			}
		})
	}
}
