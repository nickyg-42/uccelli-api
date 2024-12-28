package handlers

import (
	"encoding/json"
	"nest/db"
	"nest/models"
	"nest/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupMemberOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO Validation

	createdEvent, err := db.CreateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.DeleteGroup(r.Context(), eventID)
	if err != nil {
		http.Error(w, "Event not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAllEventsForUser(w http.ResponseWriter, r *http.Request) {
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

	events, err := db.GetAllEventsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func GetAllEventsForGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	events, err := db.GetAllEventsByGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func UpdateEventName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventName string `json:"event_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventName == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventName(r.Context(), eventID, payload.EventName)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventDescription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventDescription string `json:"event_description"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventDescription == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventDescription(r.Context(), eventID, payload.EventDescription)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventStartTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventStartTime time.Time `json:"event_start_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventStartTime(r.Context(), eventID, payload.EventStartTime)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventEndTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventEndTime time.Time `json:"event_end_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventEndTime(r.Context(), eventID, payload.EventEndTime)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
