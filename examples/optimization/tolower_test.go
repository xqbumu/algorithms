// go test -v -cpu=4 -run=none -bench=. -benchtime=10s -benchmem tolower_test.go
package optimization

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

func toLower(dst []byte, s string) []byte {
	n := len(s) - 1
	_ = s[n]
	for i := 0; i <= n; i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		dst = append(dst, c)
		fmt.Printf("address of slice %p \n", &dst)
	}
	fmt.Printf("stage address of slice %p \n", &dst)
	return dst
}

func ToLower(s string) string {
	b := make([]byte, 0, 128)
	b = toLower(b, s)
	fmt.Printf("final address of slice %p \n", &b)
	return *(*string)(unsafe.Pointer(&b))
}

func BenchmarkStdLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ToLower("The Quick Brown Fox Jumps Over The Lazy Dog")
	}
}

func BenchmarkFastLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToLower("The Quick Brown Fox Jumps Over The Lazy Dog")
	}
}

/*
BenchmarkStdLower
BenchmarkStdLower-4    	 5464696	       220 ns/op	      48 B/op	       1 allocs/op
BenchmarkFastLower
BenchmarkFastLower-4   	17009325	        69.9 ns/op	       0 B/op	       0 allocs/op
*/

func TestToLower(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"Quick Brown Fox"}, "quick brown fox"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLower(tt.args.s); got != tt.want {
				t.Errorf("ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}
