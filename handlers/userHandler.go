package handlers

import (
	"encoding/json"
	"log"
	"nest/db"
	"nest/models"
	"nest/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to access User %d's data", reqUser, id)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: Failed to find user with ID %d: %v", id, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: User %d's data successfully retrieved", id)
	json.NewEncoder(w).Encode(user)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: Failed to find user with ID %d: %v", id, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userInfo := models.UserInfo{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	log.Printf("INFO: Basic info retrieved for user %d (%s)", id, user.Username)
	json.NewEncoder(w).Encode(userInfo)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to delete User %d", reqUser, userID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to delete user %d: %v", userID, err)
		http.Error(w, "User not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d successfully deleted", userID)
	w.WriteHeader(http.StatusOK)
}

func UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update User %d's email", reqUser, id)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var emailUpdate struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&emailUpdate); err != nil {
		log.Printf("ERROR: Failed to decode email update request for user %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateUserEmail(r.Context(), id, emailUpdate.Email)
	if err != nil {
		log.Printf("ERROR: Failed to update email for user %d: %v", id, err)
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Email successfully updated for user %d to %s", id, emailUpdate.Email)
	w.WriteHeader(http.StatusOK)
}

func UpdateUserFirstName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update User %d's first name", reqUser, id)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var firstNameUpdate struct {
		FirstName string `json:"first_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&firstNameUpdate); err != nil {
		log.Printf("ERROR: Failed to decode first name update request for user %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateUserFirstName(r.Context(), id, firstNameUpdate.FirstName)
	if err != nil {
		log.Printf("ERROR: Failed to update first name for user %d: %v", id, err)
		http.Error(w, "Failed to update first name", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: First name successfully updated for user %d to %s", id, firstNameUpdate.FirstName)
	w.WriteHeader(http.StatusOK)
}

func UpdateUserLastName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update User %d's last name", reqUser, id)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var lastNameUpdate struct {
		LastName string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&lastNameUpdate); err != nil {
		log.Printf("ERROR: Failed to decode last name update request for user %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateUserLastName(r.Context(), id, lastNameUpdate.LastName)
	if err != nil {
		log.Printf("ERROR: Failed to update last name for user %d: %v", id, err)
		http.Error(w, "Failed to update last name", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Last name successfully updated for user %d to %s", id, lastNameUpdate.LastName)
	w.WriteHeader(http.StatusOK)
}
