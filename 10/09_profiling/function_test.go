package function

import (
	"testing"
)

var result int

func BenchmarkFibonacciRecursion5(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciRecursion(5)
	}
	result = fib
}

func BenchmarkFibonacciLoop5(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciLoop(5)
	}
	result = fib
}

func BenchmarkFibonacciRecursion10(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciRecursion(10)
	}
	result = fib
}

func BenchmarkFibonacciLoop10(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciLoop(10)
	}
	result = fib
}

func BenchmarkFibonacciRecursion20(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciRecursion(20)
	}
	result = fib
}

func BenchmarkFibonacciLoop20(b *testing.B) {
	var fib int
	for i := 0; i < b.N; i++ {
		fib = FibonacciLoop(20)
	}
	result = fib
}
