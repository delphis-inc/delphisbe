ALTER TABLE participants ADD COLUMN has_joined boolean NOT NULL DEFAULT FALSE;
CREATE UNIQUE INDEX CONCURRENTLY unique_disc_id_participant_idx ON participants (discussion_id, participant_id);
ALTER TABLE participants ADD CONSTRAINT unique_disc_id_participant_idx UNIQUE USING INDEX unique_disc_id_participant_idx;