ALTER TABLE discussions ADD CONSTRAINT discussions_moderator_fk_a9400b652043 FOREIGN KEY (moderator_id) REFERENCES moderators (id) MATCH FULL;

ALTER TABLE moderators ADD CONSTRAINT moderator_user_profiles_fk_cd9d71f9263a FOREIGN KEY (user_profile_id) REFERENCES user_profiles (id) MATCH FULL;

ALTER TABLE participants ADD CONSTRAINT participants_discussion_fk_21fea11dda5c FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL;
ALTER TABLE participants ADD CONSTRAINT participants_viewers_fk_a206e1523d88 FOREIGN KEY (viewer_id) REFERENCES viewers (id) MATCH FULL;
ALTER TABLE participants ADD CONSTRAINT participants_users_fk_d71eaf684c6c FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL;

ALTER TABLE posts ADD CONSTRAINT posts_discussions_fk_c34cae6d6fc5 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL;
ALTER TABLE posts ADD CONSTRAINT posts_participants_fk_c94a4fb2438b FOREIGN KEY (participant_id) REFERENCES participants (id) MATCH FULL;
ALTER TABLE posts ADD CONSTRAINT posts_post_contents_fk_777ecc8c7969 FOREIGN KEY (post_content_id) REFERENCES post_contents (id) MATCH FULL;

ALTER TABLE user_profiles ADD CONSTRAINT user_profiles_users_fk_73e968fc7272 FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL;

ALTER TABLE viewers ADD CONSTRAINT viewers_discussion_fk_6c6c7f7b2c3c FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL;
ALTER TABLE viewers ADD CONSTRAINT viewers_posts_fk_f39a03cde0e2 FOREIGN KEY (last_viewed_post_id) REFERENCES posts (id) MATCH FULL;
ALTER TABLE viewers ADD CONSTRAINT viewers_users_fk_827b069e4901 FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL;