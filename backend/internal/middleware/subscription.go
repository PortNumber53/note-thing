package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"note-thing/backend/internal/model"
)

func RequireActiveSubscription(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserID(r.Context())
			active, err := model.UserHasActiveSubscription(r.Context(), db, userID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
				return
			}
			if !active {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusPaymentRequired)
				json.NewEncoder(w).Encode(map[string]string{
					"error":   "subscription_required",
					"message": "An active subscription is required to access this resource.",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
