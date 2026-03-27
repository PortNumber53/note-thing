package handler

import (
	"net/http"

	"database/sql"

	"note-thing/backend/internal/model"
)

type SettingsHandler struct {
	DB *sql.DB
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	settings, err := model.GetUserSettings(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}
	respondJSON(w, http.StatusOK, settings)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		DefaultNotebookID *string `json:"defaultNotebookId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid input")
		return
	}
	settings, err := model.UpsertUserSettings(r.Context(), h.DB, userID, input.DefaultNotebookID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update settings")
		return
	}
	respondJSON(w, http.StatusOK, settings)
}
