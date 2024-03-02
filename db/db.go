package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DB               *sql.DB
	ConnectionString = "user=postgres dbname=bank password='Letsdoit' sslmode=disable"
)

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", ConnectionString)
	if err != nil {
		fmt.Errorf("unable to establish database connection: %s", err)
	}
	return nil
}

func Disconnect() error {
	return DB.Close()
}

