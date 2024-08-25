package main

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// O(n^2)
func bubbleSort(s []string) {
	ns := len(s)
	if ns < 2 {
		return
	}

	for i := 0; i < ns; i++ {
		for j := 0; j < ns-i-1; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}
		}
	}
}

func TestParseIdString(t *testing.T) {
	testCases := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"70002", []string{"70002"}},
		{"70002:70002", []string{"70002"}},
		{"70002a", nil},
		{"10101_10103", []string{"10101", "10102", "10103"}},
		{"10101_10164_10199", nil},
		{"10101_10164b", nil},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("case-%d\n", index), func(t *testing.T) {
			got, err := parseIdString(tc.in)
			t.Log(tc.in, tc.want, err)
			bubbleSort(got)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want: %v, got: %v\n", tc.want, got)
			}
		})
	}
}

func TestGetServices(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	want := []*Service{
		{"10040", "1.1.1.1", "13040"},
		{"10039", "2.2.2.2", "13030"},
		{"10038", "3.3.3.3", "13020"},
	}
	rows := sqlmock.NewRows([]string{"server_id", "entry_ip", "entry_port"})
	for _, v := range want {
		rows = rows.AddRow(v.Id, v.Ip, v.Port)
	}

	statement := "select sid, eip, eport from server where sid in"
	mock.ExpectQuery(statement).
		WithArgs("10040", "10039", "10038").
		WillReturnRows(rows)

	got, err := getServices(db, statement, "10040", "10039", "10038")
	if err != nil {
		t.Errorf("failed to get services: %v\n", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v, got: %v\n", want, got)
	}

	// ensure all expectations have been met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestGetServicesQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	statement := "select sid, eip, eport from server where sid in"
	prompt := fmt.Errorf("query failed")
	mock.ExpectQuery(statement).
		WithArgs("10040", "10039", "10038").
		WillReturnError(prompt)

	got, err := getServices(db, statement, "10040", "10039", "10038")
	if got != nil {
		t.Errorf("want: nil, got: %v\n", got)
	}
	if err == nil {
		t.Errorf("want: %v, got: nil\n", prompt)
	}

	// ensure all expectations have been met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestGetServicesRowFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	statement := "select sid, eip, eport from server where sid in"
	prompt := fmt.Errorf("row error")
	rows := sqlmock.NewRows([]string{"server_id", "entry_ip", "entry_port"}).
		AddRow("1", "1.1.1.1", "80").
		RowError(0, prompt)
	mock.ExpectQuery(statement).
		WithArgs("10040", "10039", "10038").
		WillReturnRows(rows)

	got, err := getServices(db, statement, "10040", "10039", "10038")
	if got != nil {
		t.Errorf("want: nil, got: %v\n", got)
	}
	if err == nil {
		t.Errorf("want: %v, got: nil\n", prompt)
	}

	// ensure all expectations have been met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}
