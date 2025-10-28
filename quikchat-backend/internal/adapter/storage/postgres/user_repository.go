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
	query := `INSERT INTO users (username, password_hash, display_name, profile_image_url) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Username, user.PasswordHash, user.DisplayName, user.ProfileImageURL).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password_hash, display_name, profile_image_url, created_at, updated_at, last_seen_at 
              FROM users WHERE username = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.ProfileImageURL,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeenAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, password_hash, display_name, profile_image_url, created_at, updated_at, last_seen_at 
              FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.ProfileImageURL,
		&user.CreatedAt, &user.UpdatedAt, &user.LastSeenAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET display_name = $1, profile_image_url = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, user.DisplayName, user.ProfileImageURL, user.ID)
	return err
}

