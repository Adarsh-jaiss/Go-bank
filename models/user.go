package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uint64         `json:"user_id"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_time"`
	Email             string         `json:"email"`
	Password          string         `json:"password"`
	PasswordChangedAt time.Time      `json:"password_changed_at"`
	CreatedAt         time.Time      `json:"created_at"`
	Accounts          []*Account     `json:"accounts"`     // Slice of accounts
	Transactions      []*Transaction `json:"transactions"` // Slice of transactions
}

type Session struct {
	ID            uuid.UUID `json:"id"`
	AccountNumber uint64    `json:"account_number"`
	Email         string    `json:"email"`
	RefreshToken  string    `json:"refresh_token"`
	UserAgent     string    `json:"user_agent"`
	ClientIP      string    `json:"client_ip"`
	IsBlocked     bool      `json:"is_blocked"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}
