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
		// User
		r.Get("/user/{id}", handlers.GetUser)

		r.Delete("/user/{id}", handlers.DeleteUser)

		r.Patch("/user/{id}/email", handlers.UpdateUserEmail)
		r.Patch("/user/{id}/firstname", handlers.UpdateUserFirstName)
		r.Patch("/user/{id}/lastname", handlers.UpdateUserLastName)

		// Group
		r.Get("/group/{id}", handlers.GetGroup)

		r.Post("/group", handlers.CreateGroup)
		r.Post("group/{id}/{user_id}", handlers.AddUserToGroup)

		r.Delete("/group/{id}", handlers.DeleteGroup)

		// Endpoint accessible only to group admins
		//r.With(middleware.RoleMiddleware("group_admin")).Get("/admin/dashboard", handlers.AdminDashboard)

		// Endpoint accessible only to super admins
		//r.With(middleware.RoleMiddleware("superadmin")).Get("/superadmin/settings", handlers.SuperAdminSettings)
	})

	return r
}
