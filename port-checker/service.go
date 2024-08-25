package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Service struct {
	// server id
	Id string
	// entry ip
	Ip string
	// entry port
	Port string
}

// getServices fetch a series of game services from MySQL server using a given SQL query and arguments.
func getServices(db *sql.DB, query string, args ...any) ([]*Service, error) {
	var sl []*Service
	rows, err := db.Query(query, args...)
	if err != nil {
		return sl, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &Service{}
		if err := rows.Scan(&s.Id, &s.Ip, &s.Port); err != nil {
			return nil, err
		}
		sl = append(sl, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sl, nil
}

func isNumericString(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// parseIdString parses an ID string of the format '10101_10164:80011_80018:90002:70002' into a slice.
func parseIdString(ids string) ([]string, error) {
	var sl []string
	if ids == "" {
		return sl, fmt.Errorf("empty string")
	}

	for _, v := range strings.Split(ids, ":") {
		if isNumericString(v) {
			sl = append(sl, v)
		} else if strings.Contains(v, "_") {
			subSlice := strings.Split(v, "_")
			for i, w := range subSlice {
				if i > 1 {
					return nil, fmt.Errorf("%s not a pair", v)
				}
				if !isNumericString(w) {
					return nil, fmt.Errorf("parsing %s: invalid syntax", w)
				}
			}
			head, _ := strconv.Atoi(subSlice[0])
			tail, _ := strconv.Atoi(subSlice[1])
			for i := head; i <= tail; i++ {
				sl = append(sl, strconv.Itoa(i))
			}
		} else {
			return nil, fmt.Errorf("parsing %s: invalid syntax", v)
		}
	}

	m := make(map[string]struct{}, len(sl))
	for _, y := range sl {
		m[y] = struct{}{}
	}
	set := make([]string, 0, len(m))
	for k := range m {
		set = append(set, k)
	}

	return set, nil
}
