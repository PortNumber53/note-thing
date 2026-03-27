ALTER TABLE notes
    DROP COLUMN IF EXISTS encrypted_title,
    DROP COLUMN IF EXISTS encrypted_body,
    DROP COLUMN IF EXISTS note_key_wrapped,
    DROP COLUMN IF EXISTS key_version,
    DROP COLUMN IF EXISTS is_encrypted;

DROP TABLE IF EXISTS user_encryption;

CREATE OR REPLACE FUNCTION notes_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
