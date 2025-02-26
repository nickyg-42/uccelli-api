package handlers

import (
	"encoding/json"
	"log"
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
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupMemberOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to access Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		log.Printf("ERROR: Failed to find event with ID %d: %v", eventID, err)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Event %d successfully retrieved by user %d", eventID, r.Context().Value("user_id").(int))
	json.NewEncoder(w).Encode(event)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var eventDTO models.EventDTO

	if err := json.NewDecoder(r.Body).Decode(&eventDTO); err != nil {
		log.Printf("ERROR: Failed to decode event creation request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNewEvent(eventDTO); err != nil {
		log.Printf("ERROR: Event validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.IsGroupMemberOrSA(r, int(eventDTO.GroupID)) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to create event in group %d", reqUser, eventDTO.GroupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	event := models.Event{
		GroupID:     eventDTO.GroupID,
		Name:        eventDTO.Name,
		Description: eventDTO.Description,
		StartTime:   eventDTO.StartTime,
		EndTime:     eventDTO.EndTime,
		CreatedByID: eventDTO.CreatedByID,
	}

	createdEvent, err := db.CreateEvent(r.Context(), &event)
	if err != nil {
		log.Printf("ERROR: Failed to create event in group %d: %v", eventDTO.GroupID, err)
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: New event created - ID: %d, Name: %s, Group: %d, Creator: %d",
		createdEvent.ID, createdEvent.Name, createdEvent.GroupID, createdEvent.CreatedByID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to delete Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Get event details before deletion for logging
	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		log.Printf("ERROR: Failed to find event %d before deletion: %v", eventID, err)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	err = db.DeleteEvent(r.Context(), eventID)
	if err != nil {
		log.Printf("ERROR: Failed to delete event %d: %v", eventID, err)
		http.Error(w, "Event not found or could not be deleted", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Event deleted - ID: %d, Name: %s, Group: %d",
		event.ID, event.Name, event.GroupID)
	w.WriteHeader(http.StatusOK)
}

func GetAllEventsForUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid user ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to access events for User %d", reqUser, userID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	events, err := db.GetAllEventsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve events for user %d: %v", userID, err)
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successfully retrieved events for user %d", userID)
	json.NewEncoder(w).Encode(events)
}

func GetAllEventsForGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid group ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Group %d not found: %v", groupID, err)
		http.Error(w, "Group does not exist", http.StatusNotFound)
		return
	}

	if !utils.IsGroupMemberOrSA(r, groupID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to access events for Group %d", reqUser, groupID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	events, err := db.GetAllEventsByGroup(r.Context(), groupID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve events for group %d: %v", groupID, err)
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successfully retrieved events for group %d", groupID)
	json.NewEncoder(w).Encode(events)
}

func UpdateEventName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventName string `json:"event_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventName == "" {
		log.Printf("ERROR: Invalid request payload for event name update: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventName(r.Context(), eventID, payload.EventName)
	if err != nil {
		log.Printf("ERROR: Failed to update event name for event %d: %v", eventID, err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Event name updated for event %d", eventID)
	w.WriteHeader(http.StatusOK)
}

func UpdateEventDescription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventDescription string `json:"event_description"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventDescription == "" {
		log.Printf("ERROR: Invalid request payload for event description update: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventDescription(r.Context(), eventID, payload.EventDescription)
	if err != nil {
		log.Printf("ERROR: Failed to update event description for event %d: %v", eventID, err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Event description updated for event %d", eventID)
	w.WriteHeader(http.StatusOK)
}

func UpdateEventStartTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventStartTime time.Time `json:"event_start_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventStartTime.IsZero() {
		log.Printf("ERROR: Invalid request payload for event start time update: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventStartTime(r.Context(), eventID, payload.EventStartTime)
	if err != nil {
		log.Printf("ERROR: Failed to update event start time for event %d: %v", eventID, err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Event start time updated for event %d", eventID)
	w.WriteHeader(http.StatusOK)
}

func UpdateEventEndTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid event ID format: %s: %v", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		EventEndTime time.Time `json:"event_end_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventEndTime.IsZero() {
		log.Printf("ERROR: Invalid request payload for event end time update: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		reqUser := r.Context().Value("user_id").(int)
		log.Printf("ERROR: Access denied - User %d attempted to update Event %d", reqUser, eventID)
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	err = db.UpdateEventEndTime(r.Context(), eventID, payload.EventEndTime)
	if err != nil {
		log.Printf("ERROR: Failed to update event end time for event %d: %v", eventID, err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Event end time updated for event %d", eventID)
	w.WriteHeader(http.StatusOK)
}
