package middleware

import (
	"database/sql"
	"net/http"
)

// RequireActiveSubscription is a pass-through middleware.
// Free tier limits (50 notes, 1 notebook, 1MB note size) are enforced
// at the handler level, not via middleware. All users can access all
// endpoints. Paid users get unlimited usage.
func RequireActiveSubscription(_ *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
