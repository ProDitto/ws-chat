package postgres

import (
	"context"
	"quikchat/internal/domain"
	"quikchat/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) repository.MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Save(ctx context.Context, message *domain.Message) error {
	query := `INSERT INTO messages (conversation_id, sender_id, type, content) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, message.ConversationID, message.SenderID, message.Type, message.Content).Scan(&message.ID, &message.CreatedAt)
	return err
}

func (r *MessageRepository) FindMessagesByConversationID(ctx context.Context, conversationID int64, cursor int64, limit int) ([]*domain.Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.type, m.content, m.created_at,
		       u.username, u.display_name, u.profile_image_url
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1 AND ($2 = 0 OR m.id < $2)
		ORDER BY m.id DESC
		LIMIT $3`

	rows, err := r.db.Query(ctx, query, conversationID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var msg domain.Message
		var sender domain.User
		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Type, &msg.Content, &msg.CreatedAt,
			&sender.Username, &sender.DisplayName, &sender.ProfileImageURL,
		)
		if err != nil {
			return nil, err
		}
		sender.ID = msg.SenderID
		msg.Sender = &sender
		messages = append(messages, &msg)
	}

	// Reverse slice to return messages in ascending order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, rows.Err()
}

