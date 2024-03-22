package store

import (
	"database/sql"
	"strconv"
	"strings"

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
	UpdateUser(*gin.Context, int64) (*models.UserAccount, error)
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

func (p *PostgresUserStore) GetUserByAccountNumber(c *gin.Context, accountNumber uint64) (*models.UserAccount, error) {
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


// UpdateUser updates user details in the PostgreSQL store
func (p *PostgresUserStore) UpdateUser(c *gin.Context, accountNumber int64) (*models.UserAccount, error) {
    // Fetch user details from request body or any other source
    var updateUser models.UserAccount
    if err := c.BindJSON(&updateUser); err != nil {
        return nil, fmt.Errorf("error binding JSON: %v", err)
    }

    // Ensure at least one field is provided for update
    if updateUser.FirstName == "" && updateUser.LastName == "" && updateUser.Email == "" {
        return nil, fmt.Errorf("no fields to update")
    }

    // Query the database to get the current user details
    var user models.UserAccount
    if updateUser.FirstName != "" || updateUser.LastName != "" || updateUser.Email != "" {
        err := p.db.QueryRow("SELECT first_name, last_name, email, balance FROM users WHERE account_number = $1", accountNumber).
            Scan(&user.FirstName, &user.LastName, &user.Email, &user.Balance)
        if err != nil {
            return nil, fmt.Errorf("error querying user details from the database: %v", err)
        }
    }

    // Generate the SQL query dynamically based on the fields provided in the request
    var queryArgs []interface{}
    var updateFields []string
    query := "UPDATE users SET"

    if updateUser.FirstName != "" && updateUser.FirstName != user.FirstName {
        updateFields = append(updateFields, "first_name = $"+strconv.Itoa(len(queryArgs)+1))
        queryArgs = append(queryArgs, updateUser.FirstName)
    }
    if updateUser.LastName != "" && updateUser.LastName != user.LastName {
        updateFields = append(updateFields, "last_name = $"+strconv.Itoa(len(queryArgs)+1))
        queryArgs = append(queryArgs, updateUser.LastName)
    }
    if updateUser.Email != "" && updateUser.Email != user.Email {
        updateFields = append(updateFields, "email = $"+strconv.Itoa(len(queryArgs)+1))
        queryArgs = append(queryArgs, updateUser.Email)
    }

    if len(updateFields) == 0 {
        return nil, fmt.Errorf("no fields to update")
    }

    query += " " + strings.Join(updateFields, ", ") + " WHERE account_number = $" + strconv.Itoa(len(queryArgs)+1) + " RETURNING first_name, last_name, email, account_number, balance"

    // Execute the SQL query
    row := p.db.QueryRow(query, append(queryArgs, accountNumber)...)
    if err := row.Scan(&updateUser.FirstName, &updateUser.LastName, &updateUser.Email, &updateUser.AccountNumber, &updateUser.Balance); err != nil {
        return nil, fmt.Errorf("error updating the user in the database: %v", err)
    }

    return &updateUser, nil
}

