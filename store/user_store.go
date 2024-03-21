package store

import (
	"database/sql"

	"fmt"

	"github.com/adarsh-jaiss/go-bank/models"
	"github.com/gin-gonic/gin"
)

// UserStorer defines the interface for user operations.
type UserStorer interface {
	InsertUser(*gin.Context, *models.User) (*models.UserAccount, error)
	GetUser(*gin.Context, *models.User) (*models.UserAccount, error)
	GetUserByAccountNumber(*gin.Context, uint64) (*models.UserAccount, error)
	GetUserByEmail(*gin.Context, string) (*models.User, error)
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

	err := p.db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password).Scan(&Createduser.FirstName, &Createduser.LastName, &Createduser.Email, &Createduser.AccountNumber, &Createduser.Balance)

	if err != nil {
		return nil, fmt.Errorf("error inserting the user in the database: %v", err)
	}

	return &Createduser, nil
}
// GetUser retrieves a user from the database based on email and password.
func (p *PostgresUserStore) GetUser(c *gin.Context, user *models.User) (*models.UserAccount, error) {
  query := `SELECT first_name, last_name, email, account_number, balance FROM users WHERE email = $1 AND password = $2`

  var userAccount models.UserAccount
  err := p.db.QueryRow(query, user.Email, user.Password).Scan(&userAccount.FirstName, &userAccount.LastName, &userAccount.Email, &userAccount.AccountNumber, &userAccount.Balance)
  if err != nil {
    return nil, fmt.Errorf("error getting the user from the database: %v", err)
  }

  return &userAccount, nil
}

func(p *PostgresUserStore) GetUserByAccountNumber(c *gin.Context, accountNumber uint64) (*models.UserAccount, error) {
	query := `SELECT first_name, last_name, email, account_number, balance FROM users WHERE account_number = $1`

	var userAccount models.UserAccount
	err := p.db.QueryRow(query, accountNumber).Scan(&userAccount.FirstName, &userAccount.LastName, &userAccount.Email, &userAccount.AccountNumber, &userAccount.Balance)
	if err != nil {
	return nil, fmt.Errorf("error getting the user from the database: %v", err)
	}

	return &userAccount, nil
	
}

func (p *PostgresUserStore) GetUserByEmail(c *gin.Context, email string) (*models.User, error) {
    query := `SELECT first_name, last_name, email, id, password FROM users WHERE email = $1`

    var user models.User
    err := p.db.QueryRow(query, email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Password)
    if err != nil {
        return nil, fmt.Errorf("error getting the user from the database: %v", err)
    }

    return &user, nil
}
