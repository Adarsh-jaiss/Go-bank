package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/adarsh-jaiss/go-bank/models"
	_ "github.com/lib/pq"
	// "github.com/adarsh-jaiss/go-bank/db"
)

var DB *sql.DB
// Seeder defines an interface for seeding functions
type Seeder interface {
    Seed(ctx context.Context, db *sql.DB) error
}

// SeedValues provides functions to seed various data types
func SeedValues(ctx context.Context, seeders []Seeder) error {
    for _, seeder := range seeders {
        err := seeder.Seed(ctx,DB)
        if err != nil {
            return fmt.Errorf("error seeding data: %w", err)
        }
        log.Printf("Seeded data for %T\n", seeder)
    }
    return nil
}

// Example seeder (replace with your specific data)
type UserSeeder struct {
    users []models.User
}

func (s *UserSeeder) Seed(ctx context.Context, db *sql.DB) error {
    for _, user := range s.users {
        // Insert user into the database
        // Replace this with your actual insert logic using prepared statements
        query := `
            INSERT INTO users (first_name, last_name, email, password)
            VALUES ($1, $2, $3, $4);
        `
        _, err := db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password)
        if err != nil {
            return fmt.Errorf("error inserting user: %w", err)
        }
    }
    return nil
}

// Example usage (replace with your actual seeders)
func main() {
    // Establish database connection (replace with your logic)
    db, err := sql.Open("postgres", "user=postgres dbname=bank password='Letsdoit' sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	
    fmt.Println("Database connected successfully!!")
    if err := Createtable(); err != nil {
		log.Fatal("Error creating table schema", err)
	}

    // Define your seeders here
    seeders := []Seeder{
        &UserSeeder{
            users: []models.User{
                {FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Password: "password"},
                // Add more users here
            },
        },
        // Add other seeders for different models
    }

    // Seed the data
    err = SeedValues(context.Background(), seeders)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Seeding completed successfully")
}



func Createtable() error {
	createUserTableSQL := `
  CREATE TABLE IF NOT EXISTS users (
      id SERIAL PRIMARY KEY,
      first_name VARCHAR(50) NOT NULL,
      last_name VARCHAR(50) NOT NULL,
      email VARCHAR(100) UNIQUE NOT NULL,
      password VARCHAR(100) NOT NULL,
	  account_number BIGINT UNIQUE NOT NULL,
      password_changed_at TIMESTAMP NOT NULL,
      created_at TIMESTAMP NOT NULL
  );`

	createUserAccountsTableSQL := `
  CREATE TABLE IF NOT EXISTS accounts (
      user_id SERIAL,
      account_number BIGINT UNIQUE NOT NULL REFERENCES users(account_number),
      balance DECIMAL(10,2) NOT NULL DEFAULT '0.00',
      FOREIGN KEY (user_id) REFERENCES users(id),
      PRIMARY KEY (user_id, account_number)
  );`

	createLedgerTableSQL := `
  CREATE TABLE IF NOT EXISTS ledger (
      id SERIAL PRIMARY KEY,
      from_account_number BIGINT NOT NULL,
      to_account_number BIGINT NOT NULL,
      type TEXT NOT NULL,
      amount DECIMAL(10,2) NOT NULL,
      created_at TIMESTAMP NOT NULL,
      FOREIGN KEY (from_account_number) REFERENCES accounts(account_number),
      FOREIGN KEY (to_account_number) REFERENCES accounts(account_number)
  );`

	createSessionTableSQL := `
  CREATE TABLE IF NOT EXISTS sessions (
      id UUID PRIMARY KEY,
      account_number BIGINT REFERENCES accounts(account_number) ON DELETE CASCADE,
      email VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE,
      refresh_token TEXT NOT NULL,
      user_agent TEXT,
      client_ip TEXT,
      is_blocked BOOLEAN,
      expires_at TIMESTAMP,
      created_at TIMESTAMP,
      FOREIGN KEY (account_number) REFERENCES accounts(account_number) ON DELETE CASCADE
  );`

	// Drop the trigger if it already exists
	dropTriggerSQL := `
      DROP TRIGGER IF EXISTS generate_account_number_trigger ON users;
      `

	// Create the trigger
	createTriggerSQL := `
  -- Create sequence for account_number with desired format
  CREATE SEQUENCE IF NOT EXISTS account_number_seq START WITH 100000000000001 INCREMENT BY 1;
  
  -- Create the trigger (function remains unchanged)
  CREATE OR REPLACE FUNCTION generate_account_number()
  RETURNS TRIGGER AS $$
  BEGIN
	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;
  
  CREATE TRIGGER generate_account_number_trigger
  BEFORE INSERT ON users
  FOR EACH ROW
  EXECUTE FUNCTION generate_account_number();
`

	// Execute SQL statements to create tables
	if _, err := DB.Exec(createUserTableSQL); err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}

	if _, err := DB.Exec(createUserAccountsTableSQL); err != nil {
		return fmt.Errorf("error creating accounts table: %v", err)
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


