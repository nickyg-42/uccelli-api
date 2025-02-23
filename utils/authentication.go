package utils

import (
	"log"
	"nest/db"
	"nest/models"
	"net/http"
)

func IsSelfOrSA(r *http.Request, userID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	if userID != authenticatedUserID && role != string(models.SuperAdmin) {
		return false
	}

	return true
}

func IsSA(r *http.Request) bool {
	role := r.Context().Value("role").(string)

	return role == string(models.SuperAdmin)
}

func IsGroupAdminOrSA(r *http.Request, groupID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	if role == string(models.SuperAdmin) {
		return true
	}

	isGroupAdmin, err := db.IsUserGroupAdmin(r.Context(), authenticatedUserID, groupID)
	if err != nil {
		return false
	}

	if role != string(models.SuperAdmin) && !isGroupAdmin {
		return false
	}

	return true
}

func IsEventCreatorOrGroupMemberOrSA(r *http.Request, eventID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		return false
	}

	isGroupMember, err := db.IsUserGroupMember(r.Context(), authenticatedUserID, int(event.GroupID))
	if err != nil {
		return false
	}

	if role != string(models.SuperAdmin) && !isGroupMember && event.CreatedByID != int64(authenticatedUserID) {
		return false
	}

	return true
}

func IsEventCreatorOrGroupAdminOrSA(r *http.Request, eventID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		return false
	}

	isGroupAdmin, err := db.IsUserGroupAdmin(r.Context(), authenticatedUserID, int(event.GroupID))
	if err != nil {
		log.Println(err)
		return false
	}

	if role != string(models.SuperAdmin) && !isGroupAdmin && event.CreatedByID != int64(authenticatedUserID) {
		return false
	}

	return true
}

func IsGroupMemberOrSA(r *http.Request, groupID int) bool {
	role := r.Context().Value("role").(string)
	authenticatedUserID := r.Context().Value("user_id").(int)

	isGroupMember, err := db.IsUserGroupMember(r.Context(), authenticatedUserID, groupID)
	if err != nil {
		return false
	}

	if role != string(models.SuperAdmin) && !isGroupMember {
		return false
	}

	return true
}
