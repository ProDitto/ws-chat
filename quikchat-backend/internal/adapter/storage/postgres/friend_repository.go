package postgres

import (
	"context"
	"errors"
	"quikchat/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendRepository struct {
	db *pgxpool.Pool
}

func NewFriendRepository(db *pgxpool.Pool) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) SaveRequest(ctx context.Context, req *domain.FriendRequest) error {
	query := `INSERT INTO friend_requests (sender_id, receiver_id, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, req.SenderID, req.ReceiverID, req.Status).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *FriendRepository) FindRequestByID(ctx context.Context, requestID int64) (*domain.FriendRequest, error) {
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at FROM friend_requests WHERE id = $1`
	req := &domain.FriendRequest{}
	err := r.db.QueryRow(ctx, query, requestID).Scan(&req.ID, &req.SenderID, &req.ReceiverID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return req, nil
}

func (r *FriendRepository) FindRequestBySenderAndReceiver(ctx context.Context, senderID, receiverID int64) (*domain.FriendRequest, error) {
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at FROM friend_requests WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)`
	req := &domain.FriendRequest{}
	err := r.db.QueryRow(ctx, query, senderID, receiverID).Scan(&req.ID, &req.SenderID, &req.ReceiverID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return req, nil
}

func (r *FriendRepository) UpdateRequestStatus(ctx context.Context, requestID int64, status domain.FriendRequestStatus) error {
	query := `UPDATE friend_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, requestID)
	return err
}

func (r *FriendRepository) FindPendingRequestsByReceiverID(ctx context.Context, userID int64) ([]*domain.FriendRequest, error) {
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at FROM friend_requests WHERE receiver_id = $1 AND status = 'pending'`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.FriendRequest
	for rows.Next() {
		req := &domain.FriendRequest{}
		if err := rows.Scan(&req.ID, &req.SenderID, &req.ReceiverID, &req.Status, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *FriendRepository) FindFriendsByUserID(ctx context.Context, userID int64) ([]*domain.User, error) {
	query := `
        SELECT u.id, u.username, u.display_name, u.profile_image_url, u.last_seen_at
        FROM users u
        JOIN friend_requests fr ON (u.id = fr.sender_id OR u.id = fr.receiver_id)
        WHERE (fr.sender_id = $1 OR fr.receiver_id = $1) AND fr.status = 'accepted' AND u.id != $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.ProfileImageURL, &user.LastSeenAt); err != nil {
			return nil, err
		}
		friends = append(friends, user)
	}
	return friends, nil
}

func (r *FriendRepository) DeleteFriendship(ctx context.Context, userID1, userID2 int64) error {
	query := `DELETE FROM friend_requests WHERE status = 'accepted' AND ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))`
	_, err := r.db.Exec(ctx, query, userID1, userID2)
	return err
}

