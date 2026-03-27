package model

import (
	"context"
	"database/sql"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func ListTags(ctx context.Context, db *sql.DB, userID string) ([]Tag, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, created_at FROM tags
		WHERE user_id = $1
		ORDER BY name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]Tag, 0)
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func CreateTag(ctx context.Context, db *sql.DB, userID, name string) (Tag, error) {
	var t Tag
	err := db.QueryRowContext(ctx, `
		INSERT INTO tags (user_id, name) VALUES ($1, $2)
		RETURNING id, name, created_at
	`, userID, name).Scan(&t.ID, &t.Name, &t.CreatedAt)
	return t, err
}

func UpdateTag(ctx context.Context, db *sql.DB, userID, tagID, name string) (Tag, error) {
	var t Tag
	err := db.QueryRowContext(ctx, `
		UPDATE tags SET name = $1 WHERE id = $2 AND user_id = $3
		RETURNING id, name, created_at
	`, name, tagID, userID).Scan(&t.ID, &t.Name, &t.CreatedAt)
	return t, err
}

func DeleteTag(ctx context.Context, db *sql.DB, userID, tagID string) error {
	result, err := db.ExecContext(ctx, `
		DELETE FROM tags WHERE id = $1 AND user_id = $2
	`, tagID, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
