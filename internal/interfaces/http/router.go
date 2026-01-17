// internal/interfaces/http/router.go
package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/invoice-app-be/internal/interfaces/http/handlers"
	mw "github.com/invoice-app-be/internal/interfaces/http/middleware"
)

type Router struct {
	invoiceHandler   *handlers.InvoiceHandler
	timeEntryHandler *handlers.TimeEntryHandler
	authHandler      *handlers.AuthHandler
	authMiddleware   *mw.AuthMiddleware
}

func NewRouter(
	invoiceHandler *handlers.InvoiceHandler,
	timeEntryHandler *handlers.TimeEntryHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *mw.AuthMiddleware,
) *Router {
	return &Router{
		invoiceHandler:   invoiceHandler,
		timeEntryHandler: timeEntryHandler,
		authHandler:      authHandler,
		authMiddleware:   authMiddleware,
	}
}

func (rt *Router) Setup() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Using chi's built-in logger
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		})

		// Auth
		r.Post("/auth/register", rt.authHandler.Register)
		r.Post("/auth/login", rt.authHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			// TODO For some reason below is now letting user in even if the user successfully logs in
			r.Use(rt.authMiddleware.Authenticate)

			// Invoices
			r.Route("/invoices", func(r chi.Router) {
				r.Get("/", rt.invoiceHandler.List)
				r.Post("/", rt.invoiceHandler.Create)
				r.Get("/{id}", rt.invoiceHandler.Get)
				r.Put("/{id}", rt.invoiceHandler.Update)
				r.Delete("/{id}", rt.invoiceHandler.Delete)
				r.Post("/{id}/send", rt.invoiceHandler.Send)
				r.Get("/{id}/pdf", rt.invoiceHandler.GeneratePDF)
			})

			// Time Entries
			r.Route("/time-entries", func(r chi.Router) {
				r.Get("/", rt.timeEntryHandler.List)
				r.Post("/", rt.timeEntryHandler.Create)
				r.Get("/{id}", rt.timeEntryHandler.Get)
				r.Put("/{id}", rt.timeEntryHandler.Update)
				r.Delete("/{id}", rt.timeEntryHandler.Delete)
				r.Post("/{id}/sync-jira", rt.timeEntryHandler.SyncToJira)
			})
		})
	})

	return r
}
