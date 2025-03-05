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
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, id) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to access Group %d", reqUser, id)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	group, err := db.GetGroupByID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: Failed to find group with ID %d: %v", id, err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Group %d successfully retrieved by user %d", id, r.Context().Value("user_id").(int))
	json.NewEncoder(w).Encode(group)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupDTO models.GroupDTO
	err := json.NewDecoder(r.Body).Decode(&groupDTO)
	if err != nil {
		log.Printf("ERROR: Failed to decode group creation request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNewGroup(groupDTO); err != nil {
		log.Printf("ERROR: Group validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := utils.GenerateRandomString(16)
	if err != nil {
		log.Printf("ERROR: Failed to generate group code: %v", err)
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
		log.Printf("ERROR: Failed to create group '%s': %v", group.Name, err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add the user as a group admin to the group they just created
	err = db.AddGroupMember(r.Context(), int(createdGroup.CreatedByID), int(createdGroup.ID), models.GroupAdmin)
	if err != nil {
		log.Printf("ERROR: Failed to add creator as admin to group %d: %v", createdGroup.ID, err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: New group created - ID: %d, Name: %s, Creator: %d",
		createdGroup.ID, createdGroup.Name, createdGroup.CreatedByID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to delete Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Get group details before deletion for logging
	group, err := db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to find group %d before deletion: %v", groupID, err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	err = db.DeleteGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to delete group %d: %v", groupID, err)
		http.Error(w, "Group not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Group deleted - ID: %d, Name: %s", group.ID, group.Name)
	w.WriteHeader(http.StatusOK)
}

func AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to add members to Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		RoleInGroup models.Role `json:"role_in_group"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.RoleInGroup == "" {
		log.Printf("ERROR: Invalid role specified for user %d in group %d: %v", userID, groupID, err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.AddGroupMember(r.Context(), userID, groupID, payload.RoleInGroup)
	if err != nil {
		log.Printf("ERROR: Failed to add user %d to group %d with role %s: %v",
			userID, groupID, payload.RoleInGroup, err)
		http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d added to group %d with role %s", userID, groupID, payload.RoleInGroup)
	w.WriteHeader(http.StatusOK)
}

func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to remove members from Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.RemoveGroupMember(r.Context(), userID, groupID)
	if err != nil {
		log.Printf("ERROR: Failed to remove user %d from group %d: %v", userID, groupID, err)
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d removed from group %d", userID, groupID)
	w.WriteHeader(http.StatusOK)
}

func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Id of the user making the request
	userID := r.Context().Value("user_id").(int)

	if !utils.IsGroupMemberOrSA(r, groupID) {
		log.Printf("ERROR: Access denied - User %d attempted to leave Group %d they're not a member of", userID, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.RemoveGroupMember(r.Context(), userID, groupID)
	if err != nil {
		log.Printf("ERROR: Failed to remove user %d from group %d: %v", userID, groupID, err)
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d left group %d", userID, groupID)
	w.WriteHeader(http.StatusOK)
}

func GetAllMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to view members of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	members, err := db.GetAllMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve members for group %d: %v", groupID, err)
		http.Error(w, "Failed to get group members", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved %d members for group %d", len(members), groupID)
	json.NewEncoder(w).Encode(members)
}

func GetAllNonMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to view non-members of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	nonMembers, err := db.GetAllNonMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve non-members for group %d: %v", groupID, err)
		http.Error(w, "Failed to get non-members", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved %d non-members for group %d", len(nonMembers), groupID)
	json.NewEncoder(w).Encode(nonMembers)
}

func GetAllNonAdminMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to view non-admin members of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	nonAdmins, err := db.GetAllNonAdminMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve non-admin members for group %d: %v", groupID, err)
		http.Error(w, "Failed to get non-admin members", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved %d non-admin members for group %d", len(nonAdmins), groupID)
	json.NewEncoder(w).Encode(nonAdmins)
}

func GetAllAdminMembersInGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to view admin members of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	admins, err := db.GetAllAdminMembersForGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve admin members for group %d: %v", groupID, err)
		http.Error(w, "Failed to get admin members", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved %d admin members for group %d", len(admins), groupID)
	json.NewEncoder(w).Encode(admins)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := db.GetAllGroups(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to retrieve all groups: %v", err)
		http.Error(w, "Failed to get groups", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved all %d groups", len(groups))
	json.NewEncoder(w).Encode(groups)
}

func GetAllGroupsForUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to view groups for User %d", reqUser, userID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	groups, err := db.GetAllGroupsForUser(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve groups for user %d: %v", userID, err)
		http.Error(w, "Failed to get groups", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Retrieved %d groups for user %d", len(groups), userID)
	json.NewEncoder(w).Encode(groups)
}

func UpdateGroupName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update name of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		GroupName string `json:"group_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.GroupName == "" {
		log.Printf("ERROR: Invalid group name update request for group %d: %v", groupID, err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateGroupName(r.Context(), groupID, payload.GroupName)
	if err != nil {
		log.Printf("ERROR: Failed to update name for group %d: %v", groupID, err)
		http.Error(w, "Failed to update group name", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Group %d name updated to '%s'", groupID, payload.GroupName)
	w.WriteHeader(http.StatusOK)
}

func UpdateGroupDoSendEmails(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update do_send_emails of Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	var payload struct {
		DoSendEmails bool `json:"do_send_emails"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("ERROR: Invalid group do_send_emails update request for group %d: %v", groupID, err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.UpdateGroupDoSendEmails(r.Context(), groupID, payload.DoSendEmails)
	if err != nil {
		log.Printf("ERROR: Failed to update do_send_emails for group %d: %v", groupID, err)
		http.Error(w, "Failed to update do_send_emails", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Group %d do_send_emails updated to '%s'", groupID, payload.DoSendEmails)
	w.WriteHeader(http.StatusOK)
}

func AddGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.AddGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to add user %d as admin to group %d: %v", userID, groupID, err)
		http.Error(w, "Failed to add group admin", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d added as admin to group %d", userID, groupID)
	w.WriteHeader(http.StatusOK)
}

func RemoveGroupAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idStr = chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.RemoveGroupAdmin(r.Context(), groupID, userID)
	if err != nil {
		log.Printf("ERROR: Failed to remove user %d as admin from group %d: %v", userID, groupID, err)
		http.Error(w, "Failed to remove group admin", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d removed as admin from group %d", userID, groupID)
	w.WriteHeader(http.StatusOK)
}

func JoinGroup(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "group_code")
	if code == "" {
		log.Printf("ERROR: Empty group code provided")
		http.Error(w, "Invalid group code", http.StatusBadRequest)
		return
	}

	group, err := db.GetGroupByCode(r.Context(), code)
	if err != nil {
		log.Printf("ERROR: Invalid group code '%s': %v", code, err)
		http.Error(w, "Invalid group code", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(int)
	err = db.AddGroupMember(r.Context(), userID, int(group.ID), models.Member)
	if err != nil {
		log.Printf("ERROR: Failed to add user %d to group %d via code: %v", userID, group.ID, err)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: User %d joined group %d using code", userID, group.ID)
	w.WriteHeader(http.StatusOK)
}
