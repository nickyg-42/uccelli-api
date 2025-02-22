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

func GetGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, id) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	group, err := db.GetGroupByID(r.Context(), id)
	if err != nil {
		log.Println("Group not found:", err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(group)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupDTO models.GroupDTO
	err := json.NewDecoder(r.Body).Decode(&groupDTO)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNewGroup(groupDTO); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := utils.GenerateRandomString(16)
	if err != nil {
		log.Println("Error generating random code:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group := models.Group{
		Name:        groupDTO.Name,
		CreatedByID: groupDTO.CreatedByID,
		Code:        code,
	}

	createdGroup, err := db.CreateGroup(r.Context(), &group)
	if err != nil {
		log.Println("Failed to create group:", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add the user as a group admin to the group they just created
	err = db.AddGroupMember(r.Context(), int(createdGroup.CreatedByID), int(createdGroup.ID), models.GroupAdmin)
	if err != nil {
		log.Println("Failed to add user to created group as an admin:", err)
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
		log.Println("Error converting ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteGroup(r.Context(), groupID)
	if err != nil {
		log.Println("Failed to delete group:", err)
		http.Error(w, "Group not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		RoleInGroup models.Role `json:"role_in_group"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.RoleInGroup == "" {
		log.Println("Error decoding request body or empty role:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.AddGroupMember(r.Context(), userID, groupID, payload.RoleInGroup)
	if err != nil {
		log.Println("Failed to add user to group:", err)
		http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.RemoveGroupMember(r.Context(), userID, groupID)
	if err != nil {
		log.Println("Failed to remove user from group:", err)
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Id of the user making the request
	userID := r.Context().Value("user_id").(int)

	if !utils.IsGroupMemberOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.RemoveGroupMember(r.Context(), userID, groupID)
	if err != nil {
		log.Println("Failed to remove user from group:", err)
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAllMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	users, err := db.GetAllMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Println("Failed to get all members in group:", err)
		http.Error(w, "Failed to get all members in group", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetAllNonMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	users, err := db.GetAllNonMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Println("Failed to get all NON members in group:", err)
		http.Error(w, "Failed to get all NON members in group", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetAllNonAdminMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	users, err := db.GetAllNonAdminMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Println("Failed to get all NON admin members in group:", err)
		http.Error(w, "Failed to get all NON admin members in group", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetAllAdminMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	users, err := db.GetAllAdminMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Println("Failed to get all admin members in group:", err)
		http.Error(w, "Failed to get all admin members in group", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if !utils.IsSA(r) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	groups, err := db.GetAllGroups(r.Context())
	if err != nil {
		log.Println("Failed to get all groups:", err)
		http.Error(w, "Failed to get all groups", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

func GetAllGroupsForUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	groups, err := db.GetAllGroupsForUser(r.Context(), userID)
	if err != nil {
		log.Println("Failed to get all groups for user:", err)
		http.Error(w, "Failed to get all groups for user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

func UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		log.Println("Access denied for resource")
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		GroupName string `json:"group_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.GroupName == "" {
		log.Println("Error decoding request body or empty group name:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateGroupName(r.Context(), groupID, payload.GroupName)
	if err != nil {
		log.Println("Failed to update group name:", err)
		http.Error(w, "Failed to update group name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AddGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.AddGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		log.Println("Failed to add group admin:", err)
		http.Error(w, "Failed to add group admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting group ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting user ID to integer:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.RemoveGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		log.Println("Failed to remove group admin:", err)
		http.Error(w, "Failed to remove group admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupCode := chi.URLParam(r, "group_code")
	userId := r.Context().Value("user_id").(int)

	group, err := db.GetGroupByCode(r.Context(), groupCode)
	if err != nil {
		log.Println("Could not find group with given code:", err)
		http.Error(w, "Invalid Code", http.StatusBadRequest)
		return
	}

	isMember, err := db.IsUserGroupMember(r.Context(), userId, int(group.ID))
	if err != nil {
		log.Println("Something went wrong checking group membership:", userId, group.ID, err)
		http.Error(w, "Something went wrong checking group membership", http.StatusBadRequest)
		return
	}

	if isMember {
		log.Println("User is already a member of group:", userId, group.ID)
		http.Error(w, "You are already a member of this group", http.StatusBadRequest)
		return
	}

	err = db.AddGroupMember(r.Context(), userId, int(group.ID), models.Member)
	if err != nil {
		log.Println("Failed to add user to group:", err)
		http.Error(w, "Failed to add user to group", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(group)
}
