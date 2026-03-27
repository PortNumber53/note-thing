package handler

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"

	"note-thing/backend/internal/model"
)

type EncryptionHandler struct {
	DB *sql.DB
}

func (h *EncryptionHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	enc, err := model.GetUserEncryption(r.Context(), h.DB, userID)
	if errors.Is(err, sql.ErrNoRows) {
		respondJSON(w, http.StatusOK, map[string]any{"enabled": false})
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get encryption metadata")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"enabled":    true,
		"kdfSalt":    enc.KDFSaltB64,
		"keyVersion": enc.KeyVersion,
		"kekVerify":  enc.KEKVerifyB64,
	})
}

func (h *EncryptionHandler) Setup(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		KDFSalt   string `json:"kdfSalt"`
		KEKVerify string `json:"kekVerify"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	salt, err := base64.StdEncoding.DecodeString(input.KDFSalt)
	if err != nil || len(salt) < 16 {
		respondError(w, http.StatusBadRequest, "invalid salt")
		return
	}
	verify, err := base64.StdEncoding.DecodeString(input.KEKVerify)
	if err != nil || len(verify) == 0 {
		respondError(w, http.StatusBadRequest, "invalid verify token")
		return
	}

	// Check if already set up
	existing, existErr := model.GetUserEncryption(r.Context(), h.DB, userID)
	if existErr == nil && existing.KeyVersion > 0 {
		respondError(w, http.StatusConflict, "encryption already set up")
		return
	}

	enc, err := model.UpsertUserEncryption(r.Context(), h.DB, userID, salt, 1, verify)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to setup encryption")
		return
	}
	respondJSON(w, http.StatusCreated, enc)
}

func (h *EncryptionHandler) RotateKey(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		KDFSalt    string `json:"kdfSalt"`
		KEKVerify  string `json:"kekVerify"`
		KeyVersion int    `json:"keyVersion"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	salt, err := base64.StdEncoding.DecodeString(input.KDFSalt)
	if err != nil || len(salt) < 16 {
		respondError(w, http.StatusBadRequest, "invalid salt")
		return
	}
	verify, err := base64.StdEncoding.DecodeString(input.KEKVerify)
	if err != nil || len(verify) == 0 {
		respondError(w, http.StatusBadRequest, "invalid verify token")
		return
	}

	enc, err := model.UpsertUserEncryption(r.Context(), h.DB, userID, salt, input.KeyVersion, verify)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to rotate key")
		return
	}
	respondJSON(w, http.StatusOK, enc)
}
