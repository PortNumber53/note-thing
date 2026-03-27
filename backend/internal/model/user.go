package model

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	GoogleID  string    `json:"-"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UpsertUser(ctx context.Context, db *sql.DB, googleID, email, name, avatarURL string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
		INSERT INTO users (google_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = now()
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, googleID, email, name, avatarURL).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func GetUserByID(ctx context.Context, db *sql.DB, id string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
		SELECT id, google_id, email, name, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func UpdateUserName(ctx context.Context, db *sql.DB, id, name string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
		UPDATE users SET name = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, name, id).Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func DeleteUser(ctx context.Context, db *sql.DB, id string) error {
	result, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
