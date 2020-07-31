BEGIN;
    ALTER TABLE discussions DROP CONSTRAINT discussions_moderator_fk_a9400b652043;
    ALTER TABLE discussions ADD CONSTRAINT discussions_moderator_fk_a9400b652043 FOREIGN KEY (moderator_id) REFERENCES moderators (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE moderators DROP CONSTRAINT moderator_user_profiles_fk_cd9d71f9263a;
    ALTER TABLE moderators ADD CONSTRAINT moderator_user_profiles_fk_cd9d71f9263a FOREIGN KEY (user_profile_id) REFERENCES user_profiles (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE participants DROP CONSTRAINT participants_discussion_fk_21fea11dda5c;
    ALTER TABLE participants ADD CONSTRAINT participants_discussion_fk_21fea11dda5c FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE participants DROP CONSTRAINT participants_viewers_fk_a206e1523d88;
    ALTER TABLE participants ADD CONSTRAINT participants_viewers_fk_a206e1523d88 FOREIGN KEY (viewer_id) REFERENCES viewers (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE participants DROP CONSTRAINT participants_users_fk_d71eaf684c6c;
    ALTER TABLE participants ADD CONSTRAINT participants_users_fk_d71eaf684c6c FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_discussions_fk_c34cae6d6fc5;
    ALTER TABLE posts ADD CONSTRAINT posts_discussions_fk_c34cae6d6fc5 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_participants_fk_c94a4fb2438b;
    ALTER TABLE posts ADD CONSTRAINT posts_participants_fk_c94a4fb2438b FOREIGN KEY (participant_id) REFERENCES participants (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_post_contents_fk_777ecc8c7969;
    ALTER TABLE posts ADD CONSTRAINT posts_post_contents_fk_777ecc8c7969 FOREIGN KEY (post_content_id) REFERENCES post_contents (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE user_profiles DROP CONSTRAINT user_profiles_users_fk_73e968fc7272;
    ALTER TABLE user_profiles ADD CONSTRAINT user_profiles_users_fk_73e968fc7272 FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE viewers DROP CONSTRAINT viewers_discussion_fk_6c6c7f7b2c3c;
    ALTER TABLE viewers ADD CONSTRAINT viewers_discussion_fk_6c6c7f7b2c3c FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE viewers DROP CONSTRAINT viewers_posts_fk_f39a03cde0e2;
    ALTER TABLE viewers ADD CONSTRAINT viewers_posts_fk_f39a03cde0e2 FOREIGN KEY (last_viewed_post_id) REFERENCES posts (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE viewers DROP CONSTRAINT viewers_users_fk_827b069e4901;
    ALTER TABLE viewers ADD CONSTRAINT viewers_users_fk_827b069e4901 FOREIGN KEY (user_id) REFERENCES users (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_quoted_post_id_eaa15cd7531b;
    ALTER TABLE posts ADD CONSTRAINT posts_quoted_post_id_eaa15cd7531b FOREIGN KEY (quoted_post_id) REFERENCES posts (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE social_infos DROP CONSTRAINT social_infos_user_profiles_fk_cc5dd3246e1d;
    ALTER TABLE social_infos ADD CONSTRAINT social_infos_user_profiles_fk_cc5dd3246e1d FOREIGN KEY (user_profile_id) REFERENCES user_profiles (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_quoted_post_id_eaa15cd7531b;
    ALTER TABLE posts ADD CONSTRAINT posts_quoted_post_id_eaa15cd7531b FOREIGN KEY (quoted_post_id) REFERENCES posts (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE activity DROP CONSTRAINT activity_participants_fk_fdc1b4ef8382;
    ALTER TABLE activity ADD CONSTRAINT activity_participants_fk_fdc1b4ef8382 FOREIGN KEY (participant_id) REFERENCES participants (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE activity DROP CONSTRAINT activity_post_contents_fk_a8bf8f3fa1b8;
    ALTER TABLE activity ADD CONSTRAINT activity_post_contents_fk_a8bf8f3fa1b8 FOREIGN KEY (post_content_id) REFERENCES post_contents (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_media_fk_b783eacfac89;
    ALTER TABLE posts ADD CONSTRAINT posts_media_fk_b783eacfac89 FOREIGN KEY (media_id) REFERENCES media (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE participants DROP CONSTRAINT participants_inviter_id_c24f6d2e8abc;
    ALTER TABLE participants ADD CONSTRAINT participants_inviter_id_c24f6d2e8abc FOREIGN KEY (inviter_id) REFERENCES participants (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE discussion_tags DROP CONSTRAINT discussions_tags_fk_abc4d1e98c22;
    ALTER TABLE discussion_tags ADD CONSTRAINT discussions_tags_fk_abc4d1e98c22 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE imported_content_tags DROP CONSTRAINT imported_content_tags_fk_0ea3db3a9fc9;
    ALTER TABLE imported_content_tags ADD CONSTRAINT imported_content_tags_fk_0ea3db3a9fc9 FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE discussion_ic_queue DROP CONSTRAINT discussion_tags_queue_fk_36fd527a4713;
    ALTER TABLE discussion_ic_queue ADD CONSTRAINT discussion_tags_queue_fk_36fd527a4713 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE discussion_ic_queue DROP CONSTRAINT imported_contents_tags_queue_fk_5222aa1ec5ca;
    ALTER TABLE discussion_ic_queue ADD CONSTRAINT imported_contents_tags_queue_fk_5222aa1ec5ca FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL ON DELETE CASCADE;

    ALTER TABLE posts DROP CONSTRAINT posts_imported_content_fk_7316d54d7c74;
    ALTER TABLE posts ADD CONSTRAINT posts_imported_content_fk_7316d54d7c74 FOREIGN KEY (imported_content_id) REFERENCES imported_contents (id) MATCH FULL ON DELETE CASCADE;
COMMIT;