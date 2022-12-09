package snkit

import (
	"strings"
)

// O I Q
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

// EncodeBit converts a base10 number to any alphabet.
func EncodeBit(n uint64, base uint64) string {
	if n == 0 {
		return string(alphabet[0])
	}
	encoded := ""
	for n > 0 {
		bs := n & 0x1f
		encoded = alphabet[bs:bs+1] + encoded
		n = n >> base
	}
	return encoded
}

// DecodeBit converts an encoded value to base10.
func DecodeBit(input string, base uint64) uint64 {
	if input == string(alphabet[0]) {
		return 0
	} else if input == "" {
		panic("Input must not be empty.")
	}

	decoded := uint64(0)
	for len(input) > 0 {
		if i := strings.Index(alphabet, input[0:1]); i >= 0 {
			decoded = decoded<<base | (uint64(i) & 0x1f)
		} else {
			panic(`invalid character`)
		}
		input = input[1:]
	}
	return decoded
}

// EncodeMod converts a base10 number to any alphabet.
func EncodeMod(n uint64, base uint64) string {
	if n == 0 {
		return string(alphabet[0])
	}
	encoded := ""
	for ; n > 0; n = n / base {
		encoded = string(alphabet[n%base]) + encoded
	}
	return encoded
}

// DecodeMod converts an encoded value to base10.
func DecodeMod(input string, base uint64) uint64 {
	if input == string(alphabet[0]) {
		return 0
	} else if input == "" {
		panic("Input must not be empty.")
	}

	decoded := uint64(0)
	for _, c := range input {
		alphabetIndex := uint64(strings.Index(alphabet, string(c)))
		decoded = base*decoded + alphabetIndex
	}
	return decoded
}
