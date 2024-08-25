package main

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// checkConn splits an ID string into a series of service IDs, attempts to connect to each service,
// and returns a series of service IDs that are not available.
func checkConn(db *sql.DB, sids string) ([]string, error) {
	sl, err := parseIdString(sids)
	if err != nil {
		return nil, err
	}
	arg := make([]any, 0, len(sl))
	for _, v := range sl {
		arg = append(arg, v)
	}
	placeholder := "select sid, eip, eport from server where sid in (%s)"
	sql := fmt.Sprintf(placeholder, strings.TrimSuffix(strings.Repeat("?,", len(sl)), ","))
	serviceList, err := getServices(db, sql, arg...)
	if err != nil {
		return nil, err
	}

	var refused []string
	for id := range telnetGroup(serviceList) {
		refused = append(refused, id)
	}

	return refused, nil
}

// telnetGroup checks if a group of services are accessible or not.
func telnetGroup(s []*Service) <-chan string {
	c := make(chan string)
	var wg sync.WaitGroup

	for _, v := range s {
		wg.Add(1)
		go func(w *Service) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(w.Ip, w.Port), 5*time.Second)
			if err != nil {
				c <- w.Id
				return
			}
			defer conn.Close()
		}(v)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}
