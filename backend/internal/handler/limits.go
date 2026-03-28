package handler

import (
	"database/sql"
	"net/http"

	"note-thing/backend/internal/model"
)

func isFreeTier(r *http.Request, db *sql.DB) bool {
	userID := userIDFromContext(r)
	active, err := model.UserHasActiveSubscription(r.Context(), db, userID)
	if err != nil {
		return true // fail closed — treat as free tier
	}
	return !active
}

func checkNoteLimit(r *http.Request, db *sql.DB, w http.ResponseWriter) bool {
	if !isFreeTier(r, db) {
		return true
	}
	userID := userIDFromContext(r)
	count, err := model.CountUserNotes(r.Context(), db, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return false
	}
	if count >= model.FreeMaxNotes {
		respondError(w, http.StatusForbidden, "Free plan limit: maximum 50 notes. Upgrade to create more.")
		return false
	}
	return true
}

func checkNotebookLimit(r *http.Request, db *sql.DB, w http.ResponseWriter) bool {
	if !isFreeTier(r, db) {
		return true
	}
	userID := userIDFromContext(r)
	count, err := model.CountUserNotebooks(r.Context(), db, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal error")
		return false
	}
	// Free users get 1 notebook (the default) + this limit allows creating 1 more = 2 total
	// But the spec says 1 notebook, so we count the default as the 1
	if count >= model.FreeMaxNotebooks+1 { // +1 for the auto-created default notebook
		respondError(w, http.StatusForbidden, "Free plan limit: maximum 1 notebook. Upgrade to create more.")
		return false
	}
	return true
}

func checkNoteSizeLimit(r *http.Request, db *sql.DB, w http.ResponseWriter, body string) bool {
	if !isFreeTier(r, db) {
		return true
	}
	if len(body) > model.FreeMaxNoteBytes {
		respondError(w, http.StatusForbidden, "Free plan limit: maximum 1MB per note. Upgrade for larger notes.")
		return false
	}
	return true
}
