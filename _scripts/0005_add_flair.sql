/* Script to create the tables and constraints for Flair.
 */
CREATE TABLE IF NOT EXISTS flairs (
    id varchar(36) PRIMARY KEY,
    display_name varchar(128),
    image_url text,
    source varchar(128) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    -- Uncomment this if you want to track when/if flair is verified
    -- verified_at timestamp with time zone
    deleted_at timestamp with time zone
);

/* Create the user/flair join table.
 *
 * This primary key is the combination of the user_id and flair_id which ensures
   that both are not null, and creates a unique index on their combination.
 * user_id has a foreign key constraint to ensure the referenced user exists.
 * If the user is deleted, delete the user's available flair also.
 * flair_id has a foreign key constraint to ensure the referenced flair exists.
 * If the flair is deleted, delete the user's available flair also.
 */
CREATE TABLE IF NOT EXISTS user_flairs (
    user_id varchar(36) REFERENCES users(id) ON DELETE CASCADE,
    flair_id varchar(36) REFERENCES flairs(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, flair_id)
);

-- Add optional flair_id column and foreign key to participants table.
ALTER TABLE participants
ADD COLUMN flair_id varchar(36)
REFERENCES flairs(id) ON DELETE SET NULL;
