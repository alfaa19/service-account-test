package models

// internal/db/models.go

import (
	"time"
)

// Account represents a bank account in the system
type Account struct {
	ID            int       `json:"id"`
	AccountNumber string    `json:"account_number"`
	Name          string    `json:"name"`
	NIK           string    `json:"nik"`
	PhoneNumber   string    `json:"phone_number"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
