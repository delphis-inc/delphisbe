ALTER TABLE discussions ADD COLUMN IF NOT EXISTS description text;
ALTER TABLE discussions ADD COLUMN IF NOT EXISTS title_history JSONB;
ALTER TABLE discussions ADD COLUMN IF NOT EXISTS description_history JSONB;