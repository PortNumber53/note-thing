ALTER TABLE notes
    ADD COLUMN user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    ADD COLUMN notebook_id UUID REFERENCES notebooks(id) ON DELETE SET NULL,
    ADD COLUMN updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD COLUMN deleted_at  TIMESTAMPTZ;

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_notebook_id ON notes(notebook_id);
CREATE INDEX idx_notes_deleted_at ON notes(deleted_at);
