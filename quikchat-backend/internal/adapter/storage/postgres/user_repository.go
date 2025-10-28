package postgres

import (
	"context"
	"errors"
	"quikchat/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, password_hash, display_name, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.Username, user.PasswordHash, user.DisplayName, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password_hash, display_name, profile_image_url, created_at, updated_at
              FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.ProfileImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, password_hash, display_name, profile_image_url, created_at, updated_at
              FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.ProfileImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}

