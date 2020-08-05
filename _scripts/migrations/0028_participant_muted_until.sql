ALTER TABLE participants
    ADD COLUMN IF NOT EXISTS muted_until timestamp with time zone;