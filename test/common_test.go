// go test -v -bench . -cover -benchmem module_name/package_path

package demo

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

func SimpleSwap(a, b *int) {
	// *a = *a + *b
	// *b = *a - *b
	// *a = *a - *b
	*a, *b = *b, *a
}

func Greeting(tip string) {
	fmt.Printf("Hello, %s\n", tip)
}

func sum(max int) int {
	total := 0
	for i := 1; i <= max; i++ {
		total += i
	}

	return total
}

func div(x, y float64) (float64, error) {
	if y == 0 {
		return 0, errors.New("被除数为 0")
	}
	return x / y, nil
}

func TestSwap(t *testing.T) {
	t.Log("Prepare to invoke TestSwap")
	x, y := 2, 3
	SimpleSwap(&x, &y)
	if x != 3 && y != 2 {
		t.Errorf("Expect x = 3, y = 2. x = %d, y = %d\n", x, y)
	}
	t.Log("Ready to teardown...")
}

func BenchmarkRandInt(b *testing.B) {
	var x, y int
	for i := 0; i < b.N; i++ {
		x, y = rand.Intn(100), rand.Intn(100)
		SimpleSwap(&x, &y)
	}
}

func ExampleGreeting() {
	Greeting("mike")
	Greeting("john")
	// Output:
	// Hello, mike
	// Hello, john
}

func BenchmarkWithDefer(b *testing.B) {
	myFunc := func() {
		defer func() {
			sum(10)
		}()
	}
	for i := 0; i < b.N; i++ {
		myFunc()
	}
}

func BenchmarkWithoutDefer(b *testing.B) {
	myFunc := func() {
		sum(10)
	}

	for i := 0; i < b.N; i++ {
		myFunc()
	}
}

// https://go.dev/wiki/TableDrivenTests
// https://go.dev/blog/subtests
func TestDivParallel(t *testing.T) {
	t.Parallel() // marks DivParallel as capable of running in parallel with other tests
	testCases := []struct {
		a, b, want float64
	}{
		{0, 0, 0},
		{0, 2, 0},
		{-1, -2, 0.5},
		{1, 1, 1},
		{2, 1, 2},
	}
	for index, tc := range testCases {
		tc := tc // capture range variable, it is necessary for Go < 1.22
		t.Run(fmt.Sprintf("Case-%d\n", index), func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			if resp, err := div(tc.a, tc.b); resp != tc.want {
				if err != nil {
					t.Log(err)
				}
				t.Errorf("Got %f, want %f\n", resp, tc.want)
			}
		})
	}
}
