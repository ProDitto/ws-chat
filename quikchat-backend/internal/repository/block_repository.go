package repository

import "context"

type BlockRepository interface {
	Save(ctx context.Context, blockerID, blockedID int64) error
	Delete(ctx context.Context, blockerID, blockedID int64) error
	IsBlocked(ctx context.Context, userID1, userID2 int64) (bool, error)
}

