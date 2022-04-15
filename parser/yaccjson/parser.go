package yaccjson

import (
	"strings"
)

func ParseJson(input string, debug bool) (interface{}, error) {
	s := NewScanner(strings.NewReader(input), debug)
	yyParse(s)
	if s.err != nil {
		return nil, s.err
	}
	return s.data, nil
}
