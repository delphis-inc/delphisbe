ALTER TABLE discussions
    ALTER COLUMN last_post_id DROP NOT NULL,
    ALTER COLUMN last_post_created_at DROP NOT NULL;