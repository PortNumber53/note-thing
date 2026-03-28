package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"note-thing/backend/internal/model"

	"github.com/go-chi/chi/v5"
)

type NotebooksHandler struct {
	DB *sql.DB
}

func (h *NotebooksHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	notebooks, err := model.ListNotebooks(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list notebooks")
		return
	}
	respondJSON(w, http.StatusOK, notebooks)
}

func (h *NotebooksHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil || input.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if !checkNotebookLimit(r, h.DB, w) {
		return
	}
	nb, err := model.CreateNotebook(r.Context(), h.DB, userID, input.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create notebook")
		return
	}
	respondJSON(w, http.StatusCreated, nb)
}

func (h *NotebooksHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	notebookID := chi.URLParam(r, "notebookID")
	var input struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil || input.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	nb, err := model.UpdateNotebook(r.Context(), h.DB, userID, notebookID, input.Name)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "notebook not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update notebook")
		return
	}
	respondJSON(w, http.StatusOK, nb)
}

func (h *NotebooksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	notebookID := chi.URLParam(r, "notebookID")
	err := model.DeleteNotebook(r.Context(), h.DB, userID, notebookID)
	if errors.Is(err, model.ErrCannotDeleteDefault) {
		respondError(w, http.StatusBadRequest, "cannot delete default notebook")
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "notebook not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete notebook")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
