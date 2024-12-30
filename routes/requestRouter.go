package routes

import (
	"nest/handlers"
	"nest/middleware"
	"nest/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", handlers.Login)
	r.Post("/register", handlers.Register)

	// JWT required routes
	r.With(middleware.JWTAuthMiddleware).Group(func(r chi.Router) {
		// User
		r.Get("/user/{id}", handlers.GetUser)
		r.Get("/user/{id}/event", handlers.GetAllEventsForUser)

		r.Delete("/user/{id}", handlers.DeleteUser)

		r.Patch("/user/{id}/email", handlers.UpdateUserEmail)
		r.Patch("/user/{id}/firstname", handlers.UpdateUserFirstName)
		r.Patch("/user/{id}/lastname", handlers.UpdateUserLastName)

		// Group
		r.Get("/group/{id}", handlers.GetGroup)
		r.Get("/group/{id}/user", handlers.GetAllMembersInGroup)
		r.Get("/group/{id}/event", handlers.GetAllEventsForGroup)

		r.Post("/group", handlers.CreateGroup)
		r.Post("group/{id}/user/{user_id}", handlers.AddUserToGroup)

		r.Patch("/group/{id}/name", handlers.UpdateGroupName)

		r.Delete("/group/{id}/user/{user_id}", handlers.RemoveUserFromGroup)
		r.Delete("/group/{id}", handlers.DeleteGroup)

		// Event
		r.Get("/event/{id}", handlers.GetEvent)

		r.Post("/event", handlers.CreateEvent)

		r.Patch("/event/{id}/name", handlers.UpdateEventName)
		r.Patch("/event/{id}/description", handlers.UpdateEventDescription)
		r.Patch("/event/{id}/start", handlers.UpdateEventStartTime)
		r.Patch("/event/{id}/end", handlers.UpdateEventEndTime)

		r.Delete("/event/{id}", handlers.DeleteEvent)

		// SA endpoints
		r.With(middleware.RoleMiddleware(models.SuperAdmin)).Patch("/group/{id}/admin/add/{user_id}", handlers.AddGroupAdmin)
		r.With(middleware.RoleMiddleware(models.SuperAdmin)).Patch("/group/{id}/admin/remove/{user_id}", handlers.RemoveGroupAdmin)
	})

	return r
}
