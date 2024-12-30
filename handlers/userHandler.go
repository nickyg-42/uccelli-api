package handlers

import (
	"encoding/json"
	"nest/db"
	"nest/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		Email string `json:"email"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateEmail(payload.Email) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserEmail(r.Context(), id, payload.Email)
	if err != nil {
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserFirstName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		FirstName string `json:"first_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateName(payload.FirstName) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserFirstName(r.Context(), id, payload.FirstName)
	if err != nil {
		http.Error(w, "Failed to update first name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserLastName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		LastName string `json:"last_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateName(payload.LastName) {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserLastName(r.Context(), id, payload.LastName)
	if err != nil {
		http.Error(w, "Failed to update last name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
