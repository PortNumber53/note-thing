package model

import (
	"context"
	"database/sql"
	"time"
)

type Notebook struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"isDefault"`
	NoteCount int       `json:"noteCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func CreateDefaultNotebook(ctx context.Context, tx *sql.Tx, userID string) (Notebook, error) {
	var nb Notebook
	err := tx.QueryRowContext(ctx, `
		INSERT INTO notebooks (user_id, name, is_default)
		VALUES ($1, 'My Notes', true)
		RETURNING id, user_id, name, is_default, created_at, updated_at
	`, userID).Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.IsDefault, &nb.CreatedAt, &nb.UpdatedAt)
	nb.NoteCount = 0
	return nb, err
}

func ListNotebooks(ctx context.Context, db *sql.DB, userID string) ([]Notebook, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT n.id, n.user_id, n.name, n.is_default, n.created_at, n.updated_at,
			   COUNT(nt.id) FILTER (WHERE nt.deleted_at IS NULL) AS note_count
		FROM notebooks n
		LEFT JOIN notes nt ON nt.notebook_id = n.id
		WHERE n.user_id = $1
		GROUP BY n.id
		ORDER BY n.is_default DESC, n.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notebooks := make([]Notebook, 0)
	for rows.Next() {
		var nb Notebook
		if err := rows.Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.IsDefault, &nb.CreatedAt, &nb.UpdatedAt, &nb.NoteCount); err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}
	return notebooks, rows.Err()
}

func CreateNotebook(ctx context.Context, db *sql.DB, userID, name string) (Notebook, error) {
	var nb Notebook
	err := db.QueryRowContext(ctx, `
		INSERT INTO notebooks (user_id, name)
		VALUES ($1, $2)
		RETURNING id, user_id, name, is_default, created_at, updated_at
	`, userID, name).Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.IsDefault, &nb.CreatedAt, &nb.UpdatedAt)
	nb.NoteCount = 0
	return nb, err
}

func UpdateNotebook(ctx context.Context, db *sql.DB, userID, notebookID, name string) (Notebook, error) {
	var nb Notebook
	err := db.QueryRowContext(ctx, `
		UPDATE notebooks SET name = $1, updated_at = now()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, name, is_default, created_at, updated_at
	`, name, notebookID, userID).Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.IsDefault, &nb.CreatedAt, &nb.UpdatedAt)
	return nb, err
}

func DeleteNotebook(ctx context.Context, db *sql.DB, userID, notebookID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get default notebook
	var defaultID string
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM notebooks WHERE user_id = $1 AND is_default = true
	`, userID).Scan(&defaultID)
	if err != nil {
		return err
	}

	if defaultID == notebookID {
		return ErrCannotDeleteDefault
	}

	// Move notes to default notebook
	_, err = tx.ExecContext(ctx, `
		UPDATE notes SET notebook_id = $1, updated_at = now()
		WHERE notebook_id = $2 AND user_id = $3
	`, defaultID, notebookID, userID)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM notebooks WHERE id = $1 AND user_id = $2
	`, notebookID, userID)
	if err != nil {
		return err
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}
