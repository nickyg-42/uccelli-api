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
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		log.Println("User not found:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(r.Context(), id)
	if err != nil {
		log.Println("User not found:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userInfo := models.UserInfo{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	json.NewEncoder(w).Encode(userInfo)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteUser(r.Context(), userID)
	if err != nil {
		log.Println("User not found or could not be deleted:", err)
		http.Error(w, "User not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		Email string `json:"email"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateEmail(payload.Email) {
		log.Println("Invalid request payload or email:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserEmail(r.Context(), id, payload.Email)
	if err != nil {
		log.Println("Failed to update email:", err)
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserFirstName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		FirstName string `json:"first_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateName(payload.FirstName) {
		log.Println("Invalid request payload or first name:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserFirstName(r.Context(), id, payload.FirstName)
	if err != nil {
		log.Println("Failed to update first name:", err)
		http.Error(w, "Failed to update first name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserLastName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, id) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		LastName string `json:"last_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || !utils.ValidateName(payload.LastName) {
		log.Println("Invalid request payload or last name:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserLastName(r.Context(), id, payload.LastName)
	if err != nil {
		log.Println("Failed to update last name:", err)
		http.Error(w, "Failed to update last name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
