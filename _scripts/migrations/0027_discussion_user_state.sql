ALTER TABLE discussion_user_access
    ADD COLUMN IF NOT EXISTS state varchar(10) DEFAULT 'ACTIVE' not null,
    ADD COLUMN IF NOT EXISTS request_id varchar(36);

ALTER TABLE discussion_user_access
    ADD CONSTRAINT dua_request_fk_ad48ef76e433 FOREIGN KEY (request_id) REFERENCES discussion_user_requests (id) MATCH FULL ON DELETE CASCADE;