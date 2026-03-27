package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"note-thing/backend/internal/billing"
	"note-thing/backend/internal/model"
)

type AdminHandler struct {
	DB      *sql.DB
	Billing *billing.Service
}

func (h *AdminHandler) ChangePrice(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AmountCents int `json:"amountCents"`
	}
	if err := decodeJSON(r, &input); err != nil || input.AmountCents <= 0 {
		respondError(w, http.StatusBadRequest, "amountCents must be positive")
		return
	}

	migration, err := h.Billing.ChangePrice(r.Context(), input.AmountCents)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to change price")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"migration": migration})
}

func (h *AdminHandler) MigrationStatus(w http.ResponseWriter, r *http.Request) {
	migration, err := model.GetLatestPriceMigration(r.Context(), h.DB)
	if errors.Is(err, sql.ErrNoRows) {
		respondJSON(w, http.StatusOK, map[string]any{"migration": nil})
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get migration")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"migration": migration})
}
