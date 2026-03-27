package model

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type Note struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	NotebookID *string   `json:"notebookId"`
	Tags       []Tag     `json:"tags"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type NoteFilters struct {
	NotebookID string
	TagID      string
	Trashed    bool
}

type CreateNoteInput struct {
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	NotebookID *string  `json:"notebookId"`
	TagIDs     []string `json:"tagIds"`
}

type UpdateNoteInput struct {
	Title      *string  `json:"title"`
	Body       *string  `json:"body"`
	NotebookID *string  `json:"notebookId"`
}

func ListNotes(ctx context.Context, db *sql.DB, userID string, filters NoteFilters) ([]Note, error) {
	query := `
		SELECT DISTINCT n.id, n.title, n.body, n.notebook_id, n.created_at, n.updated_at
		FROM notes n
	`
	args := []any{userID}
	argIdx := 2

	if filters.TagID != "" {
		query += ` JOIN note_tags nt ON nt.note_id = n.id`
	}

	query += ` WHERE n.user_id = $1`

	if filters.Trashed {
		query += ` AND n.deleted_at IS NOT NULL`
	} else {
		query += ` AND n.deleted_at IS NULL`
	}

	if filters.NotebookID != "" {
		query += ` AND n.notebook_id = $` + itoa(argIdx)
		args = append(args, filters.NotebookID)
		argIdx++
	}

	if filters.TagID != "" {
		query += ` AND nt.tag_id = $` + itoa(argIdx)
		args = append(args, filters.TagID)
		argIdx++
	}

	_ = argIdx
	query += ` ORDER BY n.updated_at DESC`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.NotebookID, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		n.Tags = make([]Tag, 0)
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		if err := loadTagsForNotes(ctx, db, notes); err != nil {
			return nil, err
		}
	}

	return notes, nil
}

func GetNote(ctx context.Context, db *sql.DB, userID, noteID string) (Note, error) {
	var n Note
	err := db.QueryRowContext(ctx, `
		SELECT id, title, body, notebook_id, created_at, updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2
	`, noteID, userID).Scan(&n.ID, &n.Title, &n.Body, &n.NotebookID, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return n, err
	}

	n.Tags = make([]Tag, 0)
	notes := []Note{n}
	if err := loadTagsForNotes(ctx, db, notes); err != nil {
		return n, err
	}
	n.Tags = notes[0].Tags
	return n, nil
}

func CreateNote(ctx context.Context, db *sql.DB, userID string, input CreateNoteInput) (Note, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return Note{}, err
	}
	defer tx.Rollback()

	// If no notebook specified, use default
	notebookID := input.NotebookID
	if notebookID == nil {
		var defaultID string
		err := tx.QueryRowContext(ctx, `
			SELECT id FROM notebooks WHERE user_id = $1 AND is_default = true
		`, userID).Scan(&defaultID)
		if err != nil {
			return Note{}, err
		}
		notebookID = &defaultID
	}

	var n Note
	err = tx.QueryRowContext(ctx, `
		INSERT INTO notes (title, body, user_id, notebook_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, body, notebook_id, created_at, updated_at
	`, input.Title, input.Body, userID, notebookID).Scan(
		&n.ID, &n.Title, &n.Body, &n.NotebookID, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return Note{}, err
	}

	n.Tags = make([]Tag, 0)
	if len(input.TagIDs) > 0 {
		if err := setNoteTags(ctx, tx, n.ID, input.TagIDs); err != nil {
			return Note{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Note{}, err
	}

	// Reload to get tags
	return GetNote(ctx, db, userID, n.ID)
}

func UpdateNote(ctx context.Context, db *sql.DB, userID, noteID string, input UpdateNoteInput) (Note, error) {
	var n Note
	err := db.QueryRowContext(ctx, `
		UPDATE notes SET
			title = COALESCE($1, title),
			body = COALESCE($2, body),
			notebook_id = COALESCE($3, notebook_id),
			updated_at = now()
		WHERE id = $4 AND user_id = $5 AND deleted_at IS NULL
		RETURNING id, title, body, notebook_id, created_at, updated_at
	`, input.Title, input.Body, input.NotebookID, noteID, userID).Scan(
		&n.ID, &n.Title, &n.Body, &n.NotebookID, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return n, err
	}

	n.Tags = make([]Tag, 0)
	notes := []Note{n}
	if err := loadTagsForNotes(ctx, db, notes); err != nil {
		return n, err
	}
	n.Tags = notes[0].Tags
	return n, nil
}

func SoftDeleteNote(ctx context.Context, db *sql.DB, userID, noteID string) error {
	result, err := db.ExecContext(ctx, `
		UPDATE notes SET deleted_at = now(), updated_at = now()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, noteID, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func RestoreNote(ctx context.Context, db *sql.DB, userID, noteID string) error {
	result, err := db.ExecContext(ctx, `
		UPDATE notes SET deleted_at = NULL, updated_at = now()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NOT NULL
	`, noteID, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func PermanentDeleteNote(ctx context.Context, db *sql.DB, userID, noteID string) error {
	result, err := db.ExecContext(ctx, `
		DELETE FROM notes WHERE id = $1 AND user_id = $2
	`, noteID, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func SearchNotes(ctx context.Context, db *sql.DB, userID, query string) ([]Note, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, title, body, notebook_id, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		  AND deleted_at IS NULL
		  AND search_vector @@ plainto_tsquery('english', $2)
		ORDER BY ts_rank(search_vector, plainto_tsquery('english', $2)) DESC
		LIMIT 50
	`, userID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.NotebookID, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		n.Tags = make([]Tag, 0)
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		if err := loadTagsForNotes(ctx, db, notes); err != nil {
			return nil, err
		}
	}

	return notes, nil
}

func SetNoteTags(ctx context.Context, db *sql.DB, userID, noteID string, tagIDs []string) error {
	// Verify note ownership
	var exists bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM notes WHERE id = $1 AND user_id = $2)`, noteID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := setNoteTags(ctx, tx, noteID, tagIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func setNoteTags(ctx context.Context, tx *sql.Tx, noteID string, tagIDs []string) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM note_tags WHERE note_id = $1`, noteID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2)
		`, noteID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadTagsForNotes(ctx context.Context, db *sql.DB, notes []Note) error {
	noteIDs := make([]string, len(notes))
	noteIndex := make(map[string]int)
	for i, n := range notes {
		noteIDs[i] = n.ID
		noteIndex[n.ID] = i
	}

	rows, err := db.QueryContext(ctx, `
		SELECT nt.note_id, t.id, t.name
		FROM note_tags nt
		JOIN tags t ON t.id = nt.tag_id
		WHERE nt.note_id = ANY($1)
		ORDER BY t.name
	`, pq.Array(noteIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var noteID, tagID, tagName string
		if err := rows.Scan(&noteID, &tagID, &tagName); err != nil {
			return err
		}
		if idx, ok := noteIndex[noteID]; ok {
			notes[idx].Tags = append(notes[idx].Tags, Tag{ID: tagID, Name: tagName})
		}
	}
	return rows.Err()
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
