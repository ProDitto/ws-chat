package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockRepository struct {
	db *pgxpool.Pool
}

func NewBlockRepository(db *pgxpool.Pool) *BlockRepository {
	return &BlockRepository{db: db}
}

func (r *BlockRepository) Save(ctx context.Context, blockerID, blockedID int64) error {
	query := `INSERT INTO user_blocks (blocker_id, blocked_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, blockerID, blockedID)
	return err
}

func (r *BlockRepository) Delete(ctx context.Context, blockerID, blockedID int64) error {
	query := `DELETE FROM user_blocks WHERE blocker_id = $1 AND blocked_id = $2`
	_, err := r.db.Exec(ctx, query, blockerID, blockedID)
	return err
}

func (r *BlockRepository) IsBlocked(ctx context.Context, userID1, userID2 int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_blocks WHERE (blocker_id = $1 AND blocked_id = $2) OR (blocker_id = $2 AND blocked_id = $1))`
	var isBlocked bool
	err := r.db.QueryRow(ctx, query, userID1, userID2).Scan(&isBlocked)
	return isBlocked, err
}

