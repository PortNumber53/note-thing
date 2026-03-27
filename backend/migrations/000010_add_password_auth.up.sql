ALTER TABLE users ADD COLUMN password_hash TEXT;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_google_id_key;
ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL;

CREATE UNIQUE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
