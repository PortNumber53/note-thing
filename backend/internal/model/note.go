package model

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"
)

const noteSelectCols = `n.id, n.title, n.body, n.encrypted_title, n.encrypted_body, n.note_key_wrapped, n.key_version, n.is_encrypted, n.notebook_id, n.created_at, n.updated_at`

type Note struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	EncryptedTitle  []byte    `json:"encryptedTitle,omitempty"`
	EncryptedBody   []byte    `json:"encryptedBody,omitempty"`
	NoteKeyWrapped  []byte    `json:"noteKeyWrapped,omitempty"`
	KeyVersion      *int      `json:"keyVersion,omitempty"`
	IsEncrypted     bool      `json:"isEncrypted"`
	NotebookID      *string   `json:"notebookId"`
	Tags            []Tag     `json:"tags"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func scanNote(row interface{ Scan(...any) error }) (Note, error) {
	var n Note
	err := row.Scan(
		&n.ID, &n.Title, &n.Body,
		&n.EncryptedTitle, &n.EncryptedBody, &n.NoteKeyWrapped, &n.KeyVersion, &n.IsEncrypted,
		&n.NotebookID, &n.CreatedAt, &n.UpdatedAt,
	)
	return n, err
}

type NoteFilters struct {
	NotebookID string
	TagID      string
	Trashed    bool
}

type CreateNoteInput struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	EncryptedTitle []byte   `json:"encryptedTitle,omitempty"`
	EncryptedBody  []byte   `json:"encryptedBody,omitempty"`
	NoteKeyWrapped []byte   `json:"noteKeyWrapped,omitempty"`
	KeyVersion     *int     `json:"keyVersion,omitempty"`
	IsEncrypted    bool     `json:"isEncrypted"`
	NotebookID     *string  `json:"notebookId"`
	TagIDs         []string `json:"tagIds"`
}

type UpdateNoteInput struct {
	Title          *string `json:"title"`
	Body           *string `json:"body"`
	EncryptedTitle []byte  `json:"encryptedTitle,omitempty"`
	EncryptedBody  []byte  `json:"encryptedBody,omitempty"`
	NoteKeyWrapped []byte  `json:"noteKeyWrapped,omitempty"`
	KeyVersion     *int    `json:"keyVersion,omitempty"`
	IsEncrypted    *bool   `json:"isEncrypted,omitempty"`
	NotebookID     *string `json:"notebookId"`
}

func ListNotes(ctx context.Context, db *sql.DB, userID string, filters NoteFilters) ([]Note, error) {
	query := `SELECT DISTINCT ` + noteSelectCols + ` FROM notes n`
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
		n, err := scanNote(rows)
		if err != nil {
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
	row := db.QueryRowContext(ctx, `
		SELECT `+noteSelectCols+`
		FROM notes n
		WHERE n.id = $1 AND n.user_id = $2
	`, noteID, userID)
	n, err := scanNote(row)
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

	notebookID := input.NotebookID
	if notebookID == nil {
		var defaultID string
		err := tx.QueryRowContext(ctx, `
			SELECT id FROM notebooks WHERE user_id = $1 AND is_default = true
		`, userID).Scan(&defaultID)
		if errors.Is(err, sql.ErrNoRows) {
			// Auto-create default notebook if missing
			nb, createErr := CreateDefaultNotebook(ctx, tx, userID)
			if createErr != nil {
				return Note{}, createErr
			}
			defaultID = nb.ID
		} else if err != nil {
			return Note{}, err
		}
		notebookID = &defaultID
	}

	var n Note
	err = tx.QueryRowContext(ctx, `
		INSERT INTO notes (title, body, encrypted_title, encrypted_body, note_key_wrapped, key_version, is_encrypted, user_id, notebook_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, body, encrypted_title, encrypted_body, note_key_wrapped, key_version, is_encrypted, notebook_id, created_at, updated_at
	`, input.Title, input.Body, input.EncryptedTitle, input.EncryptedBody, input.NoteKeyWrapped, input.KeyVersion, input.IsEncrypted, userID, notebookID,
	).Scan(
		&n.ID, &n.Title, &n.Body,
		&n.EncryptedTitle, &n.EncryptedBody, &n.NoteKeyWrapped, &n.KeyVersion, &n.IsEncrypted,
		&n.NotebookID, &n.CreatedAt, &n.UpdatedAt,
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

	return GetNote(ctx, db, userID, n.ID)
}

func UpdateNote(ctx context.Context, db *sql.DB, userID, noteID string, input UpdateNoteInput) (Note, error) {
	row := db.QueryRowContext(ctx, `
		UPDATE notes SET
			title = COALESCE($1, title),
			body = COALESCE($2, body),
			notebook_id = COALESCE($3, notebook_id),
			encrypted_title = COALESCE($4, encrypted_title),
			encrypted_body = COALESCE($5, encrypted_body),
			note_key_wrapped = COALESCE($6, note_key_wrapped),
			key_version = COALESCE($7, key_version),
			is_encrypted = COALESCE($8, is_encrypted),
			updated_at = now()
		WHERE id = $9 AND user_id = $10 AND deleted_at IS NULL
		RETURNING id, title, body, encrypted_title, encrypted_body, note_key_wrapped, key_version, is_encrypted, notebook_id, created_at, updated_at
	`, input.Title, input.Body, input.NotebookID,
		input.EncryptedTitle, input.EncryptedBody, input.NoteKeyWrapped, input.KeyVersion, input.IsEncrypted,
		noteID, userID,
	)
	n, err := scanNote(row)
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
		SELECT `+noteSelectCols+`
		FROM notes n
		WHERE n.user_id = $1
		  AND n.deleted_at IS NULL
		  AND n.is_encrypted = false
		  AND n.search_vector @@ plainto_tsquery('english', $2)
		ORDER BY ts_rank(n.search_vector, plainto_tsquery('english', $2)) DESC
		LIMIT 50
	`, userID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		n, err := scanNote(rows)
		if err != nil {
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
