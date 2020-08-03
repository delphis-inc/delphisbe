ALTER TABLE discussions
    ALTER COLUMN last_post_id SET NOT NULL,
    ALTER COLUMN last_post_created_at SET NOT NULL;