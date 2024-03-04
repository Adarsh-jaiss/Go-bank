package store

import (
	"database/sql"
	
	"fmt"

	"github.com/adarsh-jaiss/go-bank/models"
	"github.com/gin-gonic/gin"
)

// UserStorer defines the interface for user operations.
type UserStorer interface {
	InsertUser(*gin.Context, *models.User) (uint64,error)
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
func (p *PostgresUserStore) InsertUser(c *gin.Context, newUser *models.User) (uint64,error) {
    
    query := `
    INSERT INTO users (first_name, last_name, email, password, password_changed_at, created_at)
    VALUES ($1, $2, $3, $4, NOW(), NOW())
    RETURNING id;  -- Add RETURNING clause to get the inserted ID
  `
  
  var userID uint64
  
  err := p.db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password).Scan(&userID)
  
  if err != nil {
      return 0, fmt.Errorf("error inserting the user in the database: %v", err)
  }
  
//   // Marshal the user ID into JSON
//   userData, err := json.Marshal(map[string]interface{}{"id": userID})
//   if err != nil {
//       return "", fmt.Errorf("error marshalling user ID to JSON: %v", err)
//   }
  
  return userID, nil

}
