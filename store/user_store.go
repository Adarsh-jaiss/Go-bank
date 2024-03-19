package store

import (
	"database/sql"

	"fmt"

	"github.com/adarsh-jaiss/go-bank/models"
	"github.com/gin-gonic/gin"
)

// UserStorer defines the interface for user operations.
type UserStorer interface {
	InsertUser(*gin.Context, *models.User) (*models.UserAccount,error)
}

// PostgresUserStore implements UserStorer for PostgreSQL database.
type PostgresUserStore struct {
	db *sql.DB
}

// NewPostgresUserStore creates a new instance of PostgresUserStore.
func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

// InsertUser inserts a new user into the database.
func (p *PostgresUserStore) InsertUser(c *gin.Context, newUser *models.User) (*models.UserAccount, error) {
  query := `
  INSERT INTO users (first_name, last_name, email, password, password_changed_at, created_at)
  VALUES ($1, $2, $3, $4, NOW(), NOW())
  RETURNING first_name,last_name, email,account_number,balance;  -- Add RETURNING clause to get the inserted ID
`

  var Createduser models.UserAccount

  err := p.db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password).Scan(&Createduser.FirstName,&Createduser.LastName,&Createduser.Email, &Createduser.AccountNumber,&Createduser.Balance)

  if err != nil {
      return nil, fmt.Errorf("error inserting the user in the database: %v", err)
  }

  return &Createduser, nil
}

