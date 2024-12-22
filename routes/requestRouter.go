package routes

import (
	"nest/handlers"
	"nest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", handlers.Login)
	r.Post("/register", handlers.Register)

	// Authenticated Routes
	r.With(middleware.JWTAuthMiddleware).Group(func(r chi.Router) {
		// Endpoint accessible to all authenticated users
		r.Get("/user/{id}", handlers.GetUser)

		// Endpoint accessible only to admins
		//r.With(middleware.RoleMiddleware("admin")).Get("/admin/dashboard", handlers.AdminDashboard)

		// Endpoint accessible only to super admins
		//r.With(middleware.RoleMiddleware("superadmin")).Get("/superadmin/settings", handlers.SuperAdminSettings)
	})

	return r
}
