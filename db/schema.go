package db

import "fmt"

func CreateTable() error {
    createUserTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50) NOT NULL,
        last_name VARCHAR(50) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL,
        password_changed_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL
    );`

    createAccountTableSQL := `
    CREATE TABLE IF NOT EXISTS accounts (
        id SERIAL PRIMARY KEY,
        account_number BIGINT UNIQUE NOT NULL,
        user_id SERIAL REFERENCES users(id),
        balance BIGINT NOT NULL DEFAULT 0,
        currency VARCHAR(3) NOT NULL,
        created_at TIMESTAMP NOT NULL
    );`

    createUserAccountsTableSQL := `
    CREATE TABLE IF NOT EXISTS user_accounts (
        user_id SERIAL,
        account_number BIGINT,
        balance BIGINT NOT NULL DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (account_number) REFERENCES accounts(account_number),
        PRIMARY KEY (user_id, account_number)
    );`

    createLedgerTableSQL := `
    CREATE TABLE IF NOT EXISTS ledger (
        id SERIAL PRIMARY KEY,
        from_account_number BIGINT NOT NULL,
        to_account_number BIGINT NOT NULL,
        type TEXT NOT NULL,
        amount BIGINT NOT NULL,
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
    CREATE OR REPLACE FUNCTION generate_account_number()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.account_number := NEXTVAL('account_number_seq');
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER generate_account_number_trigger
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION generate_account_number();
    `

    // Execute SQL statements
    if _, err := DB.Exec(createUserTableSQL); err != nil {
        return fmt.Errorf("error creating user table: %v", err)
    }

    if _, err := DB.Exec(createAccountTableSQL); err != nil {
        return fmt.Errorf("error creating account table: %v", err)
    }

    if _, err := DB.Exec(createUserAccountsTableSQL); err != nil {
        return fmt.Errorf("error creating user_accounts table: %v", err)
    }

    if _, err := DB.Exec(createLedgerTableSQL); err != nil {
        return fmt.Errorf("error creating transaction table: %v", err)
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

    return nil
}
