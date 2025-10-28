package postgres

import (
	"context"
	"quikchat/internal/domain"
	"quikchat/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversationRepository struct {
	db *pgxpool.Pool
}

func NewConversationRepository(db *pgxpool.Pool) repository.ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Save(ctx context.Context, conversation *domain.Conversation) error {
	query := `INSERT INTO conversations (type) VALUES ($1) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, conversation.Type).Scan(&conversation.ID, &conversation.CreatedAt)
	return err
}

func (r *ConversationRepository) AddParticipant(ctx context.Context, conversationID, userID int64) error {
	query := `INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, conversationID, userID)
	return err
}

func (r *ConversationRepository) IsUserInConversation(ctx context.Context, userID, conversationID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE user_id = $1 AND conversation_id = $2)`
	err := r.db.QueryRow(ctx, query, userID, conversationID).Scan(&exists)
	return exists, err
}

func (r *ConversationRepository) FindParticipantIDs(ctx context.Context, conversationID int64) ([]int64, error) {
	query := `SELECT user_id FROM conversation_participants WHERE conversation_id = $1`
	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

