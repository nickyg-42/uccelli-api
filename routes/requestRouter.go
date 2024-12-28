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
		r.Get("/user/{id}/events", handlers.GetAllEventsForUser)

		r.Delete("/user/{id}", handlers.DeleteUser)

		r.Patch("/user/{id}/email", handlers.UpdateUserEmail)
		r.Patch("/user/{id}/firstname", handlers.UpdateUserFirstName)
		r.Patch("/user/{id}/lastname", handlers.UpdateUserLastName)

		// Group
		r.Get("/group/{id}", handlers.GetGroup)
		r.Get("/group/{id}/users", handlers.GetAllMembersInGroup)
		r.Get("/group/{id}/events", handlers.GetAllEventsForGroup)

		r.Post("/group", handlers.CreateGroup)
		r.Post("group/{id}/users/{user_id}", handlers.AddUserToGroup)

		r.Delete("/group/{id}/users/{user_id}", handlers.RemoveUserFromGroup)
		r.Delete("/group/{id}", handlers.DeleteGroup)

		// Event
		r.Get("/event/{id}", handlers.GetEvent)

		r.Post("/event", handlers.CreateEvent)

		r.Delete("/event/{id}", handlers.DeleteEvent)

		// Endpoint accessible only to group admins
		//r.With(middleware.RoleMiddleware("group_admin")).Get("/admin/dashboard", handlers.AdminDashboard)

		// stuff like get ALL <any>

		// Endpoint accessible only to super admins
		//r.With(middleware.RoleMiddleware("superadmin")).Get("/superadmin/settings", handlers.SuperAdminSettings)
	})

	return r
}
