
-- Create activity table
-- TODO: Need to think more about whether we want to track mentions and reactions in the same table
CREATE TABLE IF NOT EXISTS activity (
    participant_id varchar(36) not null,
    post_content_id varchar(36) not null,
    entity_id varchar(36) not null,
    entity_type varchar(36) not null,
    created_at timestamp with time zone default current_timestamp not null,
    PRIMARY KEY(participant_id, entity_id, created_at)
);

ALTER TABLE activity
    ADD CONSTRAINT activity_participants_fk_fdc1b4ef8382 FOREIGN KEY (participant_id) REFERENCES participants (id) MATCH FULL,
    ADD CONSTRAINT activity_post_contents_fk_a8bf8f3fa1b8 FOREIGN KEY (post_content_id) REFERENCES post_contents (id) MATCH FULL;

ALTER TABLE post_contents ADD COLUMN IF NOT EXISTS mentioned_entities varchar(50)[];