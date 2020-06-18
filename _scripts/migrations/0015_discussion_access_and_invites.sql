ALTER TABLE discussions ADD COLUMN IF NOT EXISTS public_access boolean default true not null;

CREATE TABLE IF NOT EXISTS discussion_flair_access (
    discussion_id varchar(36) not null,
    flair_template_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    PRIMARY KEY(discussion_id, flair_template_id)
);

CREATE TABLE IF NOT EXISTS discussion_user_access (
    discussion_id varchar(36) not null,
    user_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    PRIMARY KEY(user_id, discussion_id)
);

CREATE TABLE IF NOT EXISTS discussion_invite_link_access (
    discussion_id varchar(36) PRIMARY KEY,
    vip_invite_link_id varchar(36),
    invite_link_id varchar(36),
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS discussion_user_invitations (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) not null,
    discussion_id varchar(36) not null,
    invite_from_participant_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    status varchar(20) not null,
    invite_type varchar(20) not null
);

CREATE TABLE IF NOT EXISTS discussion_user_requests (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) not null,
    discussion_id varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    deleted_at timestamp with time zone,
    status varchar(20) not null
);

ALTER TABLE discussion_flair_access
    ADD CONSTRAINT dfa_discussions_fk_79064bd041e0 FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE,
    ADD CONSTRAINT dfa_flair_template_fk_98f8727dbafb FOREIGN KEY (flair_template_id) REFERENCES flair_templates(id) MATCH FULL ON DELETE CASCADE;

ALTER TABLE discussion_user_access
    ADD CONSTRAINT dua_discussion_id_fk_5a5346f2b4ef FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE,
    ADD CONSTRAINT dua_user_id_fk_980847ffbee4 FOREIGN KEY (user_id) REFERENCES users(id) MATCH FULL ON DELETE CASCADE;

ALTER TABLE discussion_invite_link_access
    ADD CONSTRAINT dila_discussion_id_fk_70cb9031f1b3 FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE;

ALTER TABLE discussion_user_invitations
    ADD CONSTRAINT dui_user_id_fk_283d28765f44 FOREIGN KEY (user_id) REFERENCES users(id) MATCH FULL ON DELETE CASCADE,
    ADD CONSTRAINT dui_discussion_id_fk_7e5c4a687f7a FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE,
    ADD CONSTRAINT dui_invite_from_participant_id_fk_258a3d978143 FOREIGN KEY (invite_from_participant_id) REFERENCES participants(id) MATCH FULL ON DELETE CASCADE;

ALTER TABLE discussion_user_requests
    ADD CONSTRAINT dur_user_id_fk_fd202a2cd1f2 FOREIGN KEY (user_id) REFERENCES users(id) MATCH FULL ON DELETE CASCADE,
    ADD CONSTRAINT dur_discussion_id_fk_04ceb371c145 FOREIGN KEY (discussion_id) REFERENCES discussions(id) MATCH FULL ON DELETE CASCADE;

ALTER TABLE discussion_flair_access
    ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone default current_timestamp not null;

ALTER TABLE discussion_user_access
    ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone default current_timestamp not null;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_invite_link_access
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_user_invitations
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_flair_access
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_user_access
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussion_user_requests
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();