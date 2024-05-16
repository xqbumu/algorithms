package main

import (
	"unsafe"
)

const alphaCaseOffset = 'a' - 'A'

func toLower(dst []byte, s string) []byte {
	n := len(s) - 1
	_ = s[n]
	for i := 0; i <= n; i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += alphaCaseOffset
		}
		dst = append(dst, c)
	}
	return dst
}

func ToLower(s string) string {
	b := make([]byte, 0, 128)
	b = toLower(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

func toUpper(dst []byte, s string) []byte {
	n := len(s) - 1
	_ = s[n]
	for i := 0; i <= n; i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			c -= alphaCaseOffset
		}
		dst = append(dst, c)
	}
	return dst
}

func ToUpper(s string) string {
	b := make([]byte, 0, 128)
	b = toUpper(b, s)
	return *(*string)(unsafe.Pointer(&b))
}
