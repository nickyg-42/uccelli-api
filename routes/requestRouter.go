package routes

import (
	"nest/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Auth endpoints
	// r.Post("/login", handlers.Login)
	// r.Post("/register", handlers.Register)

	// User endpoints
	r.Get("/user/{id}", handlers.GetUser)

	return r
}
