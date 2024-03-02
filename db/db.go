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

func CreateTable() error {
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		password_changed_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`

	createAccountTableSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		account_number INT UNIQUE NOT NULL,
		user_id INT NOT NULL,
		balance BIGINT NOT NULL,
		currency TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	createTransactionTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		from_account_number INT NOT NULL,
		to_account_number INT NOT NULL,
		type TEXT NOT NULL,
		amount BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`

	// Execute SQL statements
	if _, err := DB.Exec(createUserTableSQL); err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}

	if _, err := DB.Exec(createAccountTableSQL); err != nil {
		return fmt.Errorf("error creating account table: %v", err)
	}

	if _, err := DB.Exec(createTransactionTableSQL); err != nil {
		return fmt.Errorf("error creating transaction table: %v", err)
	}

	return nil
}
