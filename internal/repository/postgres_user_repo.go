package repository

import (
	"auth-service/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &PostgresUserRepository{db:db}
}

func (r *PostgresUserRepository) Create (ctx context.Context, username, passwordHash string) (*domain.User, error) {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash, created_at`
	
	var user domain.User

	err := r.db.QueryRow(ctx, query, username, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Gagal insert user: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query user: %w", err)
	}

	return &user, nil
}