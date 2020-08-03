ALTER TABLE discussions
    ADD COLUMN IF NOT EXISTS last_post_id varchar(36),
    ADD COLUMN IF NOT EXISTS last_post_created_at timestamp with time zone;

ALTER TABLE discussions ADD CONSTRAINT last_post_discussions_400b29940e7a FOREIGN KEY (last_post_id) REFERENCES posts (id) MATCH FULL ON DELETE CASCADE;
