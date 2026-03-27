package handler

import (
	"database/sql"
	"net/http"

	"note-thing/backend/internal/model"
)

type SearchHandler struct {
	DB *sql.DB
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	q := r.URL.Query().Get("q")
	if q == "" {
		respondJSON(w, http.StatusOK, []model.Note{})
		return
	}
	notes, err := model.SearchNotes(r.Context(), h.DB, userID, q)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "search failed")
		return
	}
	respondJSON(w, http.StatusOK, notes)
}
