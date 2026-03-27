CREATE TABLE user_encryption (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    kdf_salt    BYTEA NOT NULL,
    key_version INTEGER NOT NULL DEFAULT 1,
    kek_verify  BYTEA NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE notes
    ADD COLUMN encrypted_title  BYTEA,
    ADD COLUMN encrypted_body   BYTEA,
    ADD COLUMN note_key_wrapped BYTEA,
    ADD COLUMN key_version      INTEGER,
    ADD COLUMN is_encrypted     BOOLEAN NOT NULL DEFAULT false;

CREATE OR REPLACE FUNCTION notes_search_vector_update() RETURNS trigger AS $$
BEGIN
    IF NEW.is_encrypted THEN
        NEW.search_vector := NULL;
    ELSE
        NEW.search_vector :=
            setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(NEW.body, '')), 'B');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
