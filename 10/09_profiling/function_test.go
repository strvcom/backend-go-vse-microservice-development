package function

import (
	"testing"
)

func BenchmarkFibonacciRecursion5(b *testing.B) {
	for b.Loop() {
		FibonacciRecursion(5)
	}
}

func BenchmarkFibonacciLoop5(b *testing.B) {
	for b.Loop() {
		FibonacciLoop(5)
	}
}

func BenchmarkFibonacciRecursion10(b *testing.B) {
	for b.Loop() {
		FibonacciRecursion(10)
	}
}

func BenchmarkFibonacciLoop10(b *testing.B) {
	for b.Loop() {
		FibonacciLoop(10)
	}
}

func BenchmarkFibonacciRecursion20(b *testing.B) {
	for b.Loop() {
		FibonacciRecursion(20)
	}
}

func BenchmarkFibonacciLoop20(b *testing.B) {
	for b.Loop() {
		FibonacciLoop(20)
	}
}
