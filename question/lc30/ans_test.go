package main

import (
	"algorithms/question/data"
	"reflect"
	"testing"
)

// https://leetcode-cn.com/problems/substring-with-concatenation-of-all-words/
func findSubstring(s string, words []string) []int {
	return findSubstringWithSkip1(s, words, len(words), len(words[0]), 0)
}

func inWords(words []string, word string) int {
	for i := 0; i < len(words); i++ {
		if word == words[i] {
			return i
		}
	}
	return -1
}

func findSubstringWithSkip1(s string, words []string, wLen, step int, skipLen int) []int {
	res := []int{}
	if wLen == 0 {
		return res
	}

	forLen := len(s) - step*(wLen-skipLen)
	p := -1
	for i := 0; i <= forLen; i++ {
		p = inWords(words, s[i:i+step])
		if p == -1 {
			if skipLen == 0 {
				continue
			}
			break
		}
		if wLen-skipLen == 1 { // the last one matched
			res = append(res, i)
			continue
		}
		subRes := findSubstringWithSkip1(
			s[i+step:i+step*(wLen-skipLen)],
			append(append([]string{}, words[0:p]...), words[p+1:]...),
			wLen,
			step,
			skipLen+1,
		)
		if len(subRes) > 0 {
			res = append(res, i)
		}
	}

	return res
}

func Test_findSubstring(t *testing.T) {
	type args struct {
		s     string
		words []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{``, args{"barfoothefoobarman", []string{"foo", "bar"}}, []int{0, 9}},
		{``, args{"wordgoodgoodgoodbestword", []string{"word", "good", "best", "word"}}, []int{}},
		{``, args{"barfoofoobarthefoobarman", []string{"bar", "foo", "the"}}, []int{6, 9, 12}},
		{``, args{"aaaaaaaaaaaaaa", []string{"aa", "aa"}}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{``, args{"mississippi", []string{"is"}}, []int{1, 4}},
		{``, args{data.List30Strs["long"], data.List30SWords["long"]}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findSubstring(tt.args.s, tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}
