-- Create activity table
-- TODO: Finish
-- CREATE TABLE IF NOT EXISTS activity (
--     entity_id varchar(36) not null,
--     entity_type varchar(36) not null,
--     post_content_id
-- )

ALTER TABLE post_contents ADD COLUMN IF NOT EXISTS mentioned_entities varchar(50)[];