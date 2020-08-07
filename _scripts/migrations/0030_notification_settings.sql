ALTER TABLE discussion_user_access
    ADD COLUMN IF NOT EXISTS notif_setting varchar(10) DEFAULT 'EVERYTHING' not null;