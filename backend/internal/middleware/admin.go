package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"note-thing/backend/internal/model"
)

func RequireAdmin(db *sql.DB) func(http.Handler) http.Handler {
	adminEmails := parseAdminEmails(os.Getenv("ADMIN_EMAILS"))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := UserID(r.Context())
			user, err := model.GetUserByID(r.Context(), db, userID)
			if err != nil {
				unauthorized(w)
				return
			}
			if !adminEmails[user.Email] {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func parseAdminEmails(raw string) map[string]bool {
	m := make(map[string]bool)
	for _, e := range strings.Split(raw, ",") {
		e = strings.TrimSpace(e)
		if e != "" {
			m[e] = true
		}
	}
	return m
}
