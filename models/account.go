package models

import "time"

type Account struct {
	ID            uint64    `json:"id"`
	AccountNumber uint      `json:"account_number"`
	Balance       int64     `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}