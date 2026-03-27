DROP INDEX IF EXISTS idx_notes_deleted_at;
DROP INDEX IF EXISTS idx_notes_notebook_id;
DROP INDEX IF EXISTS idx_notes_user_id;
ALTER TABLE notes
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS notebook_id,
    DROP COLUMN IF EXISTS user_id;
