DROP INDEX IF EXISTS idx_users_google_id;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- Only re-add constraint if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'users_google_id_key'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_google_id_key UNIQUE (google_id);
    END IF;
END $$;

-- Restore NOT NULL only if all values are non-null
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE google_id IS NULL) THEN
        ALTER TABLE users ALTER COLUMN google_id SET NOT NULL;
    END IF;
END $$;
