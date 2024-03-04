package models

import "time"

type Transactions struct {
	ID                string    `json:"id"`
	FromAccountNumber uint64    `json:"from_account_number"`
	ToAccountNumber   uint64    `json:"to_account_number"`
	Type              string    `json:"type"` // "deposit", "withdrawal", "transfer", etc.
	Amount            int64     `json:"amount"`
	CreatedAt         time.Time `json:"created_at"`
}

type Deposit struct {
	ID            uint64    `json:"id"`
	AccountNumber uint64    `json:"account_number"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type Withdrawal struct {
	ID            uint64    `json:"id"`
	AccountNumber uint64    `json:"account_number"`
	Amount        int64     `json:"amount"`
	DebitedAt     time.Time `json:"debited_at"`
}