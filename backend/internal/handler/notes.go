package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"note-thing/backend/internal/model"

	"github.com/go-chi/chi/v5"
)

type NotesHandler struct {
	DB *sql.DB
}

func (h *NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	filters := model.NoteFilters{
		NotebookID: r.URL.Query().Get("notebook_id"),
		TagID:      r.URL.Query().Get("tag_id"),
	}
	notes, err := model.ListNotes(r.Context(), h.DB, userID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}
	respondJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) ListTrashed(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	notes, err := model.ListNotes(r.Context(), h.DB, userID, model.NoteFilters{Trashed: true})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list trashed notes")
		return
	}
	respondJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	note, err := model.GetNote(r.Context(), h.DB, userID, noteID)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	respondJSON(w, http.StatusOK, note)
}

func (h *NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input model.CreateNoteInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !checkNoteLimit(r, h.DB, w) {
		return
	}
	if !checkNoteSizeLimit(r, h.DB, w, input.Body) {
		return
	}
	note, err := model.CreateNote(r.Context(), h.DB, userID, input)
	if err != nil {
		log.Printf("create note failed: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create note")
		return
	}
	respondJSON(w, http.StatusCreated, note)
}

func (h *NotesHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	var input model.UpdateNoteInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Body != nil && !checkNoteSizeLimit(r, h.DB, w, *input.Body) {
		return
	}
	note, err := model.UpdateNote(r.Context(), h.DB, userID, noteID, input)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update note")
		return
	}
	respondJSON(w, http.StatusOK, note)
}

func (h *NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	err := model.SoftDeleteNote(r.Context(), h.DB, userID, noteID)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete note")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotesHandler) Restore(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	err := model.RestoreNote(r.Context(), h.DB, userID, noteID)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to restore note")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotesHandler) PermanentDelete(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	err := model.PermanentDeleteNote(r.Context(), h.DB, userID, noteID)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to permanently delete note")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotesHandler) SetTags(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	noteID := chi.URLParam(r, "noteID")
	var input struct {
		TagIDs []string `json:"tagIds"`
	}
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err := model.SetNoteTags(r.Context(), h.DB, userID, noteID, input.TagIDs)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to set tags")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
