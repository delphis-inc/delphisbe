CREATE TABLE IF NOT EXISTS discussion_shuffle_time (
    discussion_id varchar(36) PRIMARY KEY,
    shuffle_time timestamp with time zone
);

ALTER TABLE discussion_shuffle_time
    ADD CONSTRAINT dst_discussion_id_fk_4fa856c56ee0 FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS discussion_shuffle_time_idx on discussion_shuffle_time (shuffle_time);

ALTER TABLE discussions ADD COLUMN shuffle_count int default 0;