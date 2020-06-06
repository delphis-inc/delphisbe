CREATE TABLE IF NOT EXISTS imported_contents (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone default current_timestamp not null,
    content_name text not null,
    content_type text not null,
    link text not null,
    overview text not null,
    source text not null
);

CREATE TABLE IF NOT EXISTS discussion_tags (
    discussion_id varchar(36) not null,
    deleted_at timestamp with time zone,
    tag varchar(40) not null,
    created_at timestamp with time zone default current_timestamp not null,
    PRIMARY KEY(discussion_id, tag)
);

CREATE TABLE IF NOT EXISTS imported_content_tags (
    imported_content_id varchar(36) not null,
    deleted_at timestamp with time zone,
    tag varchar(40) not null,
    created_at timestamp with time zone default current_timestamp not null,
    PRIMARY KEY(imported_content_id, tag)
);

CREATE TABLE IF NOT EXISTS discussion_ic_queue (
    discussion_id varchar(36) not null,
    imported_content_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    posted_at timestamp with time zone
);

ALTER TABLE discussions
    ADD COLUMN IF NOT EXISTS auto_post boolean default true not null,
    ADD COLUMN IF NOT EXISTS idle_minutes int default 300 not null;
ALTER TABLE discussion_tags ADD CONSTRAINT discussions_tags_fk_abc4d1e98c22 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL;
ALTER TABLE imported_content_tags ADD CONSTRAINT imported_content_tags_fk_0ea3db3a9fc9 FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL;

ALTER TABLE discussion_ic_queue
    ADD CONSTRAINT discussion_tags_queue_fk_36fd527a4713 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL,
    ADD CONSTRAINT imported_contents_tags_queue_fk_5222aa1ec5ca FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL;

ALTER TABLE discussion_ic_queue ADD COLUMN IF NOT EXISTS matching_tags varchar(40)[];

ALTER TABLE posts ADD COLUMN IF NOT EXISTS imported_content_id varchar(36);
ALTER TABLE posts ADD CONSTRAINT posts_imported_content_fk_7316d54d7c74 FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_ic_queue
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();