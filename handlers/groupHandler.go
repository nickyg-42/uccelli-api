package handlers

import (
	"encoding/json"
	"nest/db"
	"nest/models"
	"nest/utils"
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

	if !utils.IsSelfOrSA(r, id) {
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

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO Validation

	createdGroup, err := db.CreateGroup(r.Context(), &group)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Group not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		RoleInGroup models.Role `json:"role_in_group"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.RoleInGroup == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.AddGroupMember(r.Context(), userID, groupID, payload.RoleInGroup)
	if err != nil {
		http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.RemoveGroupMember(r.Context(), userID, groupID)
	if err != nil {
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAllMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	users, err := db.GetAllMembersForGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		GroupName string `json:"group_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.GroupName == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateGroupName(r.Context(), groupID, payload.GroupName)
	if err != nil {
		http.Error(w, "Failed to update group name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AddGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// if !utils.IsSA(r) {
	// 	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	// 	return
	// }

	err = db.AddGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		http.Error(w, "Failed to add group admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// if !utils.IsSA(r) {
	// 	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	// 	return
	// }

	err = db.RemoveGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		http.Error(w, "Failed to remove group admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
