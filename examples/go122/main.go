package main

import (
	"fmt"
	"strings"
)

func Split(s string) func(func(int, string) bool) {
	parts := strings.Split(s, " ")
	return func(yield func(int, string) bool) {
		for i, part := range parts {
			if !yield(i, part) {
				return
			}
		}
	}
}

func main() {
	str := "Alice Bob Carole David"

	for i, x := range Split(str) {
		fmt.Println(i, x)
	}
}
