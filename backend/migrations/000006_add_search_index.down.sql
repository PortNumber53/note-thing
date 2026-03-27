DROP TRIGGER IF EXISTS trg_notes_search_vector ON notes;
DROP FUNCTION IF EXISTS notes_search_vector_update;
DROP INDEX IF EXISTS idx_notes_search;
ALTER TABLE notes DROP COLUMN IF EXISTS search_vector;
