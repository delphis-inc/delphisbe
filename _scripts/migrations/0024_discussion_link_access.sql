CREATE TABLE IF NOT EXISTS discussion_access_link (
    link_slug varchar(12) PRIMARY KEY,
    discussion_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone
);

ALTER TABLE discussion_access_link
    ADD CONSTRAINT link_discussions_fk_7f54ac52a684 FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE;

CREATE INDEX CONCURRENTLY discussion_id_f5eb7c1c82d4 ON discussion_access_link(discussion_id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_access_link
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();