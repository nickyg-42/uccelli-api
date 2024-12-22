package handlers

import (
	"encoding/json"
	"nest/db"
	"nest/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id")

	// If not self or SA, deny
	if id != authenticatedUserID && role != string(models.SuperAdmin) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	user, err := db.GetGroupByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
