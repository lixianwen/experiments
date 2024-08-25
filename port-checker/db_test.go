package main

import "testing"

// This requires a real MySQL server; using a mock doesn't make sense here.
func TestPing(t *testing.T) {
	db, err := NewDB("mysql", "localhost", "3306", "root", "123456", "ht", false)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	_, err = NewDB("mysql", "localhost", "3306", "root", "1234@5678", "ht", false)
	if err == nil {
		t.Error("want: deadline exceeded error")
	}
}
