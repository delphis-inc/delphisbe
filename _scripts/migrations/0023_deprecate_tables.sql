-- Do not run until after merging in branch
ALTER TABLE discussions DROP COLUMN IF EXISTS public_access;

DROP TABLE IF EXISTS discussion_flair_access;

DROP TABLE IF EXISTS discussion_invite_link_access;