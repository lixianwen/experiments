package main

import (
	"bytes"
	"testing"
)

func TestResizeSlice(t *testing.T) {
	keyLength := 1 << 3
	testCases := []struct {
		in, want    []byte
		description string
	}{
		{nil, make([]byte, keyLength), "nil"},
		{[]byte{65}, []byte{65, 0, 0, 0, 0, 0, 0, 0}, "with one element"},
		{[]byte{65, 42, 0x7f, 'C', 'D', 'E', 'F'}, []byte{65, 42, 127, 67, 68, 69, 70, 0}, "remain one slot"},
		{[]byte{65, 42, 127, 67, 68, 69, 70, 71}, []byte{65, 42, 127, 67, 68, 69, 70, 71}, "just fit it"},
		{[]byte{65, 42, 127, 67, 68, 69, 70, 71, 72}, []byte{65, 42, 127, 67, 68, 69, 70, 71}, "overflow one element"},
		{[]byte{65, 42, 127, 67, 68, 69, 70, 71, 72, 73, 74, 75}, []byte{65, 42, 127, 67, 68, 69, 70, 71}, "overflow multi elements"},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rev := resizeSlice(tc.in, keyLength)
			if !bytes.Equal(rev, tc.want) {
				t.Fatalf("Resize slice got %v, want %v\n", rev, tc.want)
			}
		})
	}
}

func FuzzResizeSlice(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		a = resizeSlice(a, keyLength)
		if len(a) != keyLength {
			t.Errorf("Resize slice: %v, got length: %d, want: %d\n", a, len(a), keyLength)
		}
	})
}
