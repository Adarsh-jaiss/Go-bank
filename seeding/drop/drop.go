package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)


func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=bank password='Letsdoit' sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = DropTable(db)
	if err != nil {
		fmt.Println("Error dropping tables:", err)
	} else {
		fmt.Println("Tables dropped successfully")
	}

}

func DropTable(db *sql.DB) error {
	var err error
	_, err = db.Exec("DROP SEQUENCE account_number_seq;")
	if err != nil {
		return fmt.Errorf("error dropping the tables: %w", err)
	}

	_, err = db.Exec("DROP TABLE users,accounts,ledger,sessions;")
	if err != nil {
		return fmt.Errorf("error dropping the tables: %w", err)
	}

	return nil
}
