package router

import (
	"database/sql"
	"net/http"

	"note-thing/backend/internal/handler"
	"note-thing/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(db *sql.DB, jwtSecret string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogging)

	auth := &handler.AuthHandler{DB: db, JWTSecret: jwtSecret}
	notes := &handler.NotesHandler{DB: db}
	notebooks := &handler.NotebooksHandler{DB: db}
	tags := &handler.TagsHandler{DB: db}
	search := &handler.SearchHandler{DB: db}

	// Public routes
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/auth/google/login", auth.Login)
	r.Get("/auth/google/callback", auth.Callback)
	r.Get("/callback/google/oauth", auth.Callback)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSecret))

		r.Get("/api/me", auth.Me)

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
