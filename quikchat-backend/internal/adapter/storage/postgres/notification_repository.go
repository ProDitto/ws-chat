package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) repository.NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	query := `INSERT INTO notifications (user_id, type, message, actor_id, object_id)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id, created_at, read`
	err := r.db.QueryRow(ctx, query, notification.UserID, notification.Type, notification.Message, notification.ActorID, notification.ObjectID).Scan(&notification.ID, &notification.CreatedAt, &notification.Read)
	return err
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.Notification, error) {
	query := `SELECT n.id, n.user_id, n.type, n.message, n.read, n.actor_id, n.object_id, n.created_at,
                     u.id, u.username, u.display_name, u.profile_image_url
              FROM notifications n
              LEFT JOIN users u ON n.actor_id = u.id
              WHERE n.user_id = $1
              ORDER BY n.created_at DESC
              LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		actor := &domain.User{}
		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Message, &n.Read, &n.ActorID, &n.ObjectID, &n.CreatedAt,
			&actor.ID, &actor.Username, &actor.DisplayName, &actor.ProfileImageURL,
		)
		if err != nil {
			return nil, err
		}
		if n.ActorID != nil {
			n.Actor = actor
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	query := `UPDATE notifications SET read = TRUE WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, notificationID, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET read = TRUE WHERE user_id = $1 AND read = FALSE`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *NotificationRepository) CountUnread(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = FALSE`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

