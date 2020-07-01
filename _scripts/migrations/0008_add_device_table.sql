CREATE TABLE IF NOT EXISTS user_devices (
    id varchar(36) PRIMARY KEY,
    platform VARCHAR(16) NOT NULL,
    user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL,
    last_seen timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    token VARCHAR(512)
);

CREATE INDEX user_devices_user_id_last_seen on user_devices (user_id, last_seen);