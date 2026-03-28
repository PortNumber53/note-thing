package router

import (
	"database/sql"
	"net/http"

	"note-thing/backend/internal/billing"
	"note-thing/backend/internal/handler"
	"note-thing/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(db *sql.DB, jwtSecret string, billingSvc *billing.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogging)

	auth := &handler.AuthHandler{DB: db, JWTSecret: jwtSecret}
	notes := &handler.NotesHandler{DB: db}
	notebooks := &handler.NotebooksHandler{DB: db}
	tags := &handler.TagsHandler{DB: db}
	search := &handler.SearchHandler{DB: db}
	settings := &handler.SettingsHandler{DB: db}
	billingH := &handler.BillingHandler{DB: db, Billing: billingSvc}
	admin := &handler.AdminHandler{DB: db, Billing: billingSvc}

	// Public routes
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/auth/google/login", auth.GoogleLogin)
	r.Get("/auth/google/callback", auth.Callback)
	r.Get("/callback/google/oauth", auth.Callback)
	r.Post("/auth/google/token", auth.TokenExchange)
	r.Post("/auth/signup", auth.Signup)
	r.Post("/auth/login", auth.EmailLogin)
	r.Get("/api/billing/price", billingH.GetPrice)
	r.Post("/stripe/webhook", billingH.HandleWebhook)

	// Authenticated routes (no subscription required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSecret))

		r.Get("/api/me", auth.Me)
		r.Put("/api/me", auth.UpdateMe)
		r.Delete("/api/me", auth.DeleteMe)

		r.Get("/api/settings", settings.Get)
		r.Put("/api/settings", settings.Update)

		encryption := &handler.EncryptionHandler{DB: db}
		r.Get("/api/encryption", encryption.GetMetadata)
		r.Post("/api/encryption/setup", encryption.Setup)
		r.Put("/api/encryption/rotate", encryption.RotateKey)

		r.Get("/api/billing/usage", billingH.Usage)
		r.Get("/api/billing/status", billingH.Status)
		r.Post("/api/billing/checkout", billingH.CreateCheckout)
		r.Post("/api/billing/portal", billingH.CreatePortal)
		r.Post("/api/billing/cancel", billingH.CancelSubscription)
		r.Post("/api/billing/reactivate", billingH.Reactivate)

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdmin(db))
			r.Post("/api/admin/billing/price", admin.ChangePrice)
			r.Get("/api/admin/billing/migration", admin.MigrationStatus)
		})
	})

	// Authenticated + subscription required (only enforced when billing is configured)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSecret))
		if billingSvc != nil {
			r.Use(middleware.RequireActiveSubscription(db))
		}

		r.Route("/api/notes", func(r chi.Router) {
			r.Get("/", notes.List)
			r.Post("/", notes.Create)
			r.Get("/trash", notes.ListTrashed)
			r.Get("/{noteID}", notes.Get)
			r.Put("/{noteID}", notes.Update)
			r.Delete("/{noteID}", notes.Delete)
			r.Post("/{noteID}/restore", notes.Restore)
			r.Delete("/{noteID}/permanent", notes.PermanentDelete)
			r.Put("/{noteID}/tags", notes.SetTags)
		})

		r.Route("/api/notebooks", func(r chi.Router) {
			r.Get("/", notebooks.List)
			r.Post("/", notebooks.Create)
			r.Put("/{notebookID}", notebooks.Update)
			r.Delete("/{notebookID}", notebooks.Delete)
		})

		r.Route("/api/tags", func(r chi.Router) {
			r.Get("/", tags.List)
			r.Post("/", tags.Create)
			r.Put("/{tagID}", tags.Update)
			r.Delete("/{tagID}", tags.Delete)
		})

		r.Get("/api/search", search.Search)
	})

	return r
}
