package model

import (
	"context"
	"database/sql"
)

const (
	FreeMaxNotes     = 50
	FreeMaxNotebooks = 1
	FreeMaxNoteBytes = 1 * 1024 * 1024 // 1MB
)

func CountUserNotes(ctx context.Context, db *sql.DB, userID string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM notes WHERE user_id = $1 AND deleted_at IS NULL
	`, userID).Scan(&count)
	return count, err
}

func CountUserNotebooks(ctx context.Context, db *sql.DB, userID string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM notebooks WHERE user_id = $1
	`, userID).Scan(&count)
	return count, err
}
