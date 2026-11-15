package domain

import (
	"context"
	"time"
)

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*User, error) 

	FindByUsername(ctx context.Context, username string) (*User, error)
}