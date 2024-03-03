package store

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/adarsh-jaiss/go-bank/models"
    "github.com/gin-gonic/gin"
)

// UserStorer defines the interface for user operations.
type UserStorer interface {
    InsertUser(*gin.Context, *models.User) (*models.UserAccount, error)
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
    query := `INSERT INTO users (first_name, last_name, email, password, password_changed_at, created_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`

    var userID uint64
    err := p.db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password).Scan(&userID)
    if err != nil {
        return nil, fmt.Errorf("error inserting the user in the database: %v", err)
    }

    userAccount, err := p.fetchUserAccount(userID)
    if err != nil {
        return nil, err
    }

    return userAccount, nil
}

// fetchUserAccount fetches the account details for the newly created user from the user_accounts table.
func (p *PostgresUserStore) fetchUserAccount(userID uint64) (*models.UserAccount, error) {
    accountQuery := `SELECT account_number, balance, currency
                     FROM user_accounts
                     WHERE user_id = $1`

    var accountNumber uint64
    var balance float64
    

    err := p.db.QueryRow(accountQuery, userID).Scan(&accountNumber, &balance)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user account not found for user ID %d", userID)
        }
        return nil, fmt.Errorf("error fetching user account: %v", err)
    }

    userAccount := &models.UserAccount{
        UserID:        userID,
        AccountNumber: accountNumber,
        Balance:       balance,
    }

    return userAccount, nil
}

