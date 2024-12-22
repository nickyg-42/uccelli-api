package handlers

import (
	"encoding/json"
	"nest/db"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id")

	// If not self or SA, deny
	if id != authenticatedUserID && role != "sa" {
		http.Error(w, "You do not have access to this user", http.StatusForbidden)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
