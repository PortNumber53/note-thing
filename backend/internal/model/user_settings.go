package model

import (
	"context"
	"database/sql"
	"errors"
)

type UserSettings struct {
	UserID            string  `json:"-"`
	DefaultNotebookID *string `json:"defaultNotebookId"`
}

func GetUserSettings(ctx context.Context, db *sql.DB, userID string) (UserSettings, error) {
	var s UserSettings
	err := db.QueryRowContext(ctx, `
		SELECT user_id, default_notebook_id
		FROM user_settings WHERE user_id = $1
	`, userID).Scan(&s.UserID, &s.DefaultNotebookID)
	if errors.Is(err, sql.ErrNoRows) {
		return UserSettings{UserID: userID}, nil
	}
	return s, err
}

func UpsertUserSettings(ctx context.Context, db *sql.DB, userID string, defaultNotebookID *string) (UserSettings, error) {
	var s UserSettings
	err := db.QueryRowContext(ctx, `
		INSERT INTO user_settings (user_id, default_notebook_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET
			default_notebook_id = EXCLUDED.default_notebook_id,
			updated_at = now()
		RETURNING user_id, default_notebook_id
	`, userID, defaultNotebookID).Scan(&s.UserID, &s.DefaultNotebookID)
	return s, err
}
