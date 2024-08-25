package main

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCheckConn(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	want := []*Service{
		{"10038", "localhost", "22"},
		{"10039", "localhost", "22"},
		{"10040", "localhost", "22"},
	}
	rows := sqlmock.NewRows([]string{"server_id", "entry_ip", "entry_port"})
	for _, v := range want {
		rows = rows.AddRow(v.Id, v.Ip, v.Port)
	}

	statement := "select sid, eip, eport from server where sid in"
	mock.ExpectQuery(statement).
		WithArgs("10038", "10039", "10040").
		WillReturnRows(rows)

	refused, _ := checkConn(db, "10038_10040")
	if len(refused) != 0 {
		t.Errorf("Unexpected: %v\n", refused)
	}
}

func TestCheckConnOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// change these addresses when you can connect to it
	want := []*Service{
		{"10038", "google.com", "443"},
		{"10039", "youtube.com", "443"},
		{"10040", "golang.org", "443"},
	}
	rows := sqlmock.NewRows([]string{"server_id", "entry_ip", "entry_port"})
	for _, v := range want {
		rows = rows.AddRow(v.Id, v.Ip, v.Port)
	}

	statement := "select sid, eip, eport from server where sid in"
	mock.ExpectQuery(statement).
		WithArgs("10038", "10039", "10040").
		WillReturnRows(rows)

	slowContains := func(s []string, e string) (flag bool) {
		for _, v := range s {
			if v == e {
				flag = true
				break
			}
		}

		return
	}

	refused, _ := checkConn(db, "10038_10040")
	// O(n^2)
	for _, v := range want {
		if !slowContains(refused, v.Id) {
			t.Errorf("want %s in %v\n", v.Id, refused)
		}
	}
}
