package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uint64    `json:"id"`
	AccountNumber uint      `json:"account_number"`
	Balance       int64     `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

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

type Transaction struct {
	ID                string    `json:"id"`
	FromAccountNumber uint64    `json:"from_account_number"`
	ToAccountNumber   uint64    `json:"to_account_number"`
	Type              string    `json:"type"` // "deposit", "withdrawal", "transfer", etc.
	Amount            int64     `json:"amount"`
	CreatedAt         time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Deposit struct {
	ID            uint64    `json:"id"`
	AccountNumber uint64    `json:"account_number"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}
