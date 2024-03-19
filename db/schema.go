package db

import (
  "fmt"

  // Import database driver (replace with your specific driver)
  _ "github.com/lib/pq" // Example for postgres driver
)



func CreateTable() error {
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(55) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		account_number BIGINT UNIQUE,
		balance DECIMAL(10,2) NOT NULL DEFAULT '0.00',
		password_changed_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`

	createLedgerTableSQL := `
	CREATE TABLE IF NOT EXISTS ledger (
		id SERIAL PRIMARY KEY,
		from_account_number BIGINT NOT NULL,
		to_account_number BIGINT NOT NULL,
		type TEXT NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (from_account_number) REFERENCES users(account_number),
		FOREIGN KEY (to_account_number) REFERENCES users(account_number)
	);`

	createSessionTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY,
		account_number BIGINT REFERENCES users(account_number) ON DELETE CASCADE,
		email VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE,
		refresh_token TEXT NOT NULL,
		user_agent TEXT,
		client_ip TEXT,
		is_blocked BOOLEAN,
		expires_at TIMESTAMP,
		created_at TIMESTAMP,
		FOREIGN KEY (account_number) REFERENCES users(account_number) ON DELETE CASCADE
	);`

	// Drop the trigger if it already exists
	dropTriggerSQL := `
	DROP TRIGGER IF EXISTS generate_account_number_trigger ON users;
	`

	// Create the trigger
	createTriggerSQL := `
	-- Create sequence for account_number with desired format
	CREATE SEQUENCE IF NOT EXISTS account_number_seq START WITH 100000000000001 INCREMENT BY 1;

	-- Create the trigger function to generate account numbers
	CREATE OR REPLACE FUNCTION generate_account_number()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.account_number := nextval('account_number_seq');
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- Attach the trigger to the users table
	CREATE TRIGGER generate_account_number_trigger
	BEFORE INSERT ON users
	FOR EACH ROW
	EXECUTE FUNCTION generate_account_number();
	`

	// Execute SQL statements to create tables
	if _, err := DB.Exec(createUserTableSQL); err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}

	if _, err := DB.Exec(createLedgerTableSQL); err != nil {
		return fmt.Errorf("error creating ledger table: %v", err)
	}

	if _, err := DB.Exec(createSessionTableSQL); err != nil {
		return fmt.Errorf("error creating session table: %v", err)
	}

	// Execute drop trigger SQL
	if _, err := DB.Exec(dropTriggerSQL); err != nil {
		return fmt.Errorf("error dropping trigger: %v", err)
	}

	// Execute create trigger SQL
	if _, err := DB.Exec(createTriggerSQL); err != nil {
		return fmt.Errorf("error creating trigger: %v", err)
	}

	fmt.Println("Tables created successfully!!!")
	return nil
}
