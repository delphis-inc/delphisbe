BEGIN;
    ALTER TABLE discussions
        DROP COLUMN IF EXISTS auto_post,
        DROP COLUMN IF EXISTS idle_minutes;

    ALTER TABLE participants DROP COLUMN IF EXISTS flair_id;

    -- Do we want to get rid of quoted posts?
    ALTER TABLE posts DROP COLUMN IF EXISTS imported_content_id;

    DROP TABLE IF EXISTS discussion_ic_queue;
    DROP TABLE IF EXISTS discussion_tags;
    DROP TABLE IF EXISTS discussion_user_invitations;
    DROP TABLE IF EXISTS flairs;
    DROP TABLE IF EXISTS flair_templates;
    DROP TABLE IF EXISTS imported_content_tags;
    DROP TABLE IF EXISTS imported_contents;

COMMIT;