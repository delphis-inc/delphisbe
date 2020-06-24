-- Add timestamp columns to post_contents
ALTER TABLE post_contents
    ADD COLUMN IF NOT EXISTS created_at timestamp with time zone default current_timestamp not null,
    ADD COLUMN IF NOT EXISTS updated_at timestamp with time zone default current_timestamp not null;

-- Alter tables to have a default timestamp on inserts
ALTER TABLE discussions
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE moderators
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE participants
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE posts
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE social_infos
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE user_profiles
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE users
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE viewers
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE flair_templates
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE flairs
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN updated_at SET DEFAULT current_timestamp;

ALTER TABLE user_devices
    ALTER COLUMN created_at SET DEFAULT current_timestamp,
    ALTER COLUMN last_seen SET DEFAULT current_timestamp;

-- Function for trigger
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Automatically update updated_at column
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON discussions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON moderators
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON participants
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON social_infos
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_profiles
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON viewers
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON flair_templates
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON flairs
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON post_contents
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();