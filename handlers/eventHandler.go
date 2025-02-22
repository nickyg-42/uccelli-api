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
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing event ID:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupMemberOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	event, err := db.GetEventByID(r.Context(), eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		log.Println("Error retrieving event by ID:", err)
		return
	}

	json.NewEncoder(w).Encode(event)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var eventDTO models.EventDTO

	if err := json.NewDecoder(r.Body).Decode(&eventDTO); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if err := utils.ValidateNewEvent(eventDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error validating new event:", err)
		return
	}

	if !utils.IsGroupMemberOrSA(r, int(eventDTO.GroupID)) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to group:", eventDTO.GroupID)
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
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		log.Println("Error creating event:", err)
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
		log.Println("Error parsing event ID:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	err = db.DeleteEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, "Event not found or could not be deleted", http.StatusInternalServerError)
		log.Println("Error deleting event:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAllEventsForUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing user ID:", err)
		return
	}

	if !utils.IsSelfOrSA(r, userID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to user events:", userID)
		return
	}

	events, err := db.GetAllEventsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		log.Println("Error retrieving events for user:", err)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func GetAllEventsForGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing group ID:", err)
		return
	}

	_, err = db.GetGroupByID(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Group does not exist", http.StatusNotFound)
		log.Println("Error retrieving group by ID:", err)
		return
	}

	if !utils.IsGroupMemberOrSA(r, groupID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to group events:", groupID)
		return
	}

	events, err := db.GetAllEventsByGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		log.Println("Error retrieving events for group:", err)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func UpdateEventName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing event ID:", err)
		return
	}

	var payload struct {
		EventName string `json:"event_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventName == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	err = db.UpdateEventName(r.Context(), eventID, payload.EventName)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		log.Println("Error updating event name:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventDescription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing event ID:", err)
		return
	}

	var payload struct {
		EventDescription string `json:"event_description"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventDescription == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	err = db.UpdateEventDescription(r.Context(), eventID, payload.EventDescription)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		log.Println("Error updating event description:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventStartTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing event ID:", err)
		return
	}

	var payload struct {
		EventStartTime time.Time `json:"event_start_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventStartTime.IsZero() {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	err = db.UpdateEventStartTime(r.Context(), eventID, payload.EventStartTime)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		log.Println("Error updating event start time:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateEventEndTime(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Error parsing event ID:", err)
		return
	}

	var payload struct {
		EventEndTime time.Time `json:"event_end_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.EventEndTime.IsZero() {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	if !utils.IsEventCreatorOrGroupAdminOrSA(r, eventID) {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		log.Println("Forbidden access to event:", eventID)
		return
	}

	err = db.UpdateEventEndTime(r.Context(), eventID, payload.EventEndTime)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		log.Println("Error updating event end time:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
