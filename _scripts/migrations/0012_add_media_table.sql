CREATE TABLE IF NOT EXISTS media (
    id varchar(36) not null PRIMARY KEY,
    created_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    deleted_reason_code varchar(16),
    media_type varchar(16),
    media_size json not null
);

ALTER TABLE posts ADD COLUMN IF NOT EXISTS media_id varchar(36);
ALTER TABLE posts ADD CONSTRAINT posts_media_fk_b783eacfac89 FOREIGN KEY (media_id) REFERENCES media (id) MATCH FULL;