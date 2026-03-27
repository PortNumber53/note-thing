package model

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"
)

type UserEncryption struct {
	UserID     string    `json:"userId"`
	KDFSalt    []byte    `json:"-"`
	KDFSaltB64 string    `json:"kdfSalt"`
	KeyVersion int       `json:"keyVersion"`
	KEKVerify  []byte    `json:"-"`
	KEKVerifyB64 string  `json:"kekVerify"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (e *UserEncryption) encodeBase64() {
	e.KDFSaltB64 = base64.StdEncoding.EncodeToString(e.KDFSalt)
	e.KEKVerifyB64 = base64.StdEncoding.EncodeToString(e.KEKVerify)
}

func GetUserEncryption(ctx context.Context, db *sql.DB, userID string) (UserEncryption, error) {
	var e UserEncryption
	err := db.QueryRowContext(ctx, `
		SELECT user_id, kdf_salt, key_version, kek_verify, created_at, updated_at
		FROM user_encryption WHERE user_id = $1
	`, userID).Scan(&e.UserID, &e.KDFSalt, &e.KeyVersion, &e.KEKVerify, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return e, err
	}
	e.encodeBase64()
	return e, nil
}

func UpsertUserEncryption(ctx context.Context, db *sql.DB, userID string, salt []byte, keyVersion int, kekVerify []byte) (UserEncryption, error) {
	var e UserEncryption
	err := db.QueryRowContext(ctx, `
		INSERT INTO user_encryption (user_id, kdf_salt, key_version, kek_verify)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			kdf_salt = EXCLUDED.kdf_salt,
			key_version = EXCLUDED.key_version,
			kek_verify = EXCLUDED.kek_verify,
			updated_at = now()
		RETURNING user_id, kdf_salt, key_version, kek_verify, created_at, updated_at
	`, userID, salt, keyVersion, kekVerify).Scan(
		&e.UserID, &e.KDFSalt, &e.KeyVersion, &e.KEKVerify, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return e, err
	}
	e.encodeBase64()
	return e, nil
}

func UserHasEncryption(ctx context.Context, db *sql.DB, userID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM user_encryption WHERE user_id = $1)
	`, userID).Scan(&exists)
	return exists, err
}
