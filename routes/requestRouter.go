package routes

import (
	"nest/handlers"
	"nest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Auth endpoints
	r.Post("/login", handlers.Login)
	r.Post("/register", handlers.Register)

	// User endpoints (protected by JWT)
	r.With(middleware.JWTAuthMiddleware).Get("/user/{id}", handlers.GetUser)

	return r
}
