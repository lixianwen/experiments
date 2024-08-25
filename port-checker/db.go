package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewDB(driver, host, port, user, password, dbname string, lazy bool) (*sql.DB, error) {
	db, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname))
	if err != nil {
		return nil, err
	}

	if !lazy {
		pingctx, pingcancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer pingcancel()
		if err := db.PingContext(pingctx); err != nil {
			return nil, err
		}
	}

	return db, nil
}
