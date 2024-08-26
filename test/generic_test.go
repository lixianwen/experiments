package demo

import "testing"

type addable interface {
	~int | ~string
}

func genericAdd[T addable](a, b T) T {
	return a + b
}

func add(a, b int) int {
	return a + b
}

func BenchmarkGenericAdd(b *testing.B) {
	b.ReportAllocs()
	var x, y int = 42, 96
	for i := 0; i < b.N; i++ {
		genericAdd(x, y)
	}
}

func BenchmarkRoutineAdd(b *testing.B) {
	b.ReportAllocs()
	var x, y int = 42, 96
	for i := 0; i < b.N; i++ {
		add(x, y)
	}
}
