ALTER TABLE notes ADD COLUMN search_vector tsvector;

CREATE INDEX idx_notes_search ON notes USING gin(search_vector);

CREATE OR REPLACE FUNCTION notes_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_notes_search_vector
    BEFORE INSERT OR UPDATE OF title, body ON notes
    FOR EACH ROW
    EXECUTE FUNCTION notes_search_vector_update();

UPDATE notes SET search_vector =
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(body, '')), 'B');
