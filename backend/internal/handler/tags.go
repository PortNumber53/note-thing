package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"note-thing/backend/internal/model"

	"github.com/go-chi/chi/v5"
)

type TagsHandler struct {
	DB *sql.DB
}

func (h *TagsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	tags, err := model.ListTags(r.Context(), h.DB, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list tags")
		return
	}
	respondJSON(w, http.StatusOK, tags)
}

func (h *TagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	var input struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil || input.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	tag, err := model.CreateTag(r.Context(), h.DB, userID, input.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create tag")
		return
	}
	respondJSON(w, http.StatusCreated, tag)
}

func (h *TagsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	tagID := chi.URLParam(r, "tagID")
	var input struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &input); err != nil || input.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}
	tag, err := model.UpdateTag(r.Context(), h.DB, userID, tagID, input.Name)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "tag not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update tag")
		return
	}
	respondJSON(w, http.StatusOK, tag)
}

func (h *TagsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	tagID := chi.URLParam(r, "tagID")
	err := model.DeleteTag(r.Context(), h.DB, userID, tagID)
	if errors.Is(err, sql.ErrNoRows) {
		respondError(w, http.StatusNotFound, "tag not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete tag")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
