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

func scanUser(row interface{ Scan(...any) error }) (User, error) {
	var u User
	var googleID sql.NullString
	err := row.Scan(&u.ID, &googleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	u.GoogleID = googleID.String
	return u, err
}

func UpsertUser(ctx context.Context, db *sql.DB, googleID, email, name, avatarURL string) (User, error) {
	// First try: insert or update by google_id
	u, err := scanUser(db.QueryRowContext(ctx, `
		INSERT INTO users (google_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = now()
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, googleID, email, name, avatarURL))
	if err == nil {
		return u, nil
	}

	// If email already exists (from email/password signup), merge by linking google_id
	return scanUser(db.QueryRowContext(ctx, `
		UPDATE users SET
			google_id = $1,
			name = COALESCE(NULLIF(name, ''), $3),
			avatar_url = CASE WHEN avatar_url = '' THEN $4 ELSE avatar_url END,
			updated_at = now()
		WHERE email = $2
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, googleID, email, name, avatarURL))
}

func GetUserByID(ctx context.Context, db *sql.DB, id string) (User, error) {
	return scanUser(db.QueryRowContext(ctx, `
		SELECT id, google_id, email, name, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`, id))
}

func UpdateUserName(ctx context.Context, db *sql.DB, id, name string) (User, error) {
	return scanUser(db.QueryRowContext(ctx, `
		UPDATE users SET name = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, name, id))
}

func CreateUserWithPassword(ctx context.Context, db *sql.DB, email, name, passwordHash string) (User, error) {
	return scanUser(db.QueryRowContext(ctx, `
		INSERT INTO users (email, name, password_hash, avatar_url)
		VALUES ($1, $2, $3, '')
		RETURNING id, google_id, email, name, avatar_url, created_at, updated_at
	`, email, name, passwordHash))
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (User, string, error) {
	var u User
	var googleID sql.NullString
	var passwordHash sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT id, google_id, email, name, avatar_url, created_at, updated_at, password_hash
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &googleID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt, &passwordHash,
	)
	u.GoogleID = googleID.String
	return u, passwordHash.String, err
}

func SetPasswordHash(ctx context.Context, db *sql.DB, id, passwordHash string) error {
	_, err := db.ExecContext(ctx, `UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2`, passwordHash, id)
	return err
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
