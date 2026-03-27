package handler

import (
	"net/http"

	"note-thing/backend/internal/middleware"
)

func userIDFromContext(r *http.Request) string {
	return middleware.UserID(r.Context())
}
