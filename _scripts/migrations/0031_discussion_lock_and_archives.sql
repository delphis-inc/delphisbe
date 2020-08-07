ALTER TABLE discussions
    ADD COLUMN IF NOT EXISTS lock_status boolean default false not null;

CREATE TABLE IF NOT EXISTS discussion_archives (
    discussion_id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone default current_timestamp not null,
    archived jsonb not null
);

ALTER TABLE discussion_archives
    ADD CONSTRAINT archives_discussions_fk_8f15ca8d43fb FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE;
