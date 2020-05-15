ALTER TABLE posts ADD COLUMN IF NOT EXISTS quoted_post_id varchar(36);
ALTER TABLE posts ADD CONSTRAINT posts_quoted_post_id_eaa15cd7531b FOREIGN KEY (quoted_post_id) REFERENCES posts (id) MATCH FULL;