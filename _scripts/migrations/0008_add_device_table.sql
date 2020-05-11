CREATE TABLE IF NOT EXISTS device_push_tokens (
    id varchar(36) PRIMARY KEY,
    token VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS user_devices (
    id varchar(36) PRIMARY KEY,
    platform VARCHAR(16),
    push_token_id VARCHAR(36) REFERENCES device_push_tokens(id) ON DELETE CASCADE,
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL,
    last_seen timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE INDEX user_devices_user_id_last_seen on user_devices (user_id, last_seen);
CREATE INDEX user_devices_push_token_id on user_devices (push_token_id);