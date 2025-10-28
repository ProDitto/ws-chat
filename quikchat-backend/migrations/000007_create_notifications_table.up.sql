CREATE TYPE notification_type AS ENUM (
    'friend_request_received',
    'friend_request_accepted',
    'group_invite',
    'new_message'
);

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    message TEXT NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    actor_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    object_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id_created_at ON notifications (user_id, created_at DESC);

