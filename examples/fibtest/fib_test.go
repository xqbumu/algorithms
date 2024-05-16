package fibtest_test

import (
	"fmt"
	"testing"
)

func TestFibonacci(t *testing.T) {
	f := fibonacci(20)
	if f != 6765 {
		t.Errorf("fibonacci (20) = %d; want 6765", f)
	}
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func BenchmarkFibonacci(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fibonacci(20)
	}
}

func ExampleFibonacci() {
	fmt.Println(fibonacci(20)) // Output: 6765
}
