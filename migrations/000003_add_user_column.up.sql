ALTER TABLE user ADD COLUMN IF NOT EXISTS auth JSON NOT NULL;
ALTER TABLE user ADD COLUMN IF NOT EXISTS meta JSON NOT NULL;
