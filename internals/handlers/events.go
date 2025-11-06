package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/models"
)

func isValidSlotStatus(s string) bool {
    switch models.SlotStatus(s) {
    case models.SlotBusy, models.SlotSwappable, models.SlotSwapPending:
        return true
    default:
        return false
    }
}
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	default:
		http.Error(w, "Invalid user ID type", http.StatusInternalServerError)
		return
	}
	type EventInput struct {
		Title     string `json:"title"`
		StartTime string `json:"startTime"` // ISO string
		EndTime   string `json:"endTime"`
		Status    string `json:"status"`
	}

	var input EventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//status must be one of the SlotStatus values
	var status models.SlotStatus = models.SlotBusy
	  if input.Status != "" {
        if !isValidSlotStatus(input.Status) {
            http.Error(w, "Invalid status value", http.StatusBadRequest)
            return
        }
        status = models.SlotStatus(input.Status)
    }
	start, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		http.Error(w, "Invalid startTime format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		http.Error(w, "Invalid endTime format", http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "endTime must be after startTime", http.StatusBadRequest)
		return
	}

	event := models.Event{
		Title:     input.Title,
		StartTime: start,
		EndTime:   end,
		Status:    status, 
		UserID:    uid,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var events []models.Event
	if err := database.DB.Where("user_id = ?", userID).Find(&events).Error; err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 🔹 Extract event ID from URL
	  path := r.URL.Path
    // Assuming path format is /events/{id}
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        http.Error(w, "Invalid URL path", http.StatusBadRequest)
        return
    }
    idStr := parts[len(parts)-1]
    
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

	type UpdateInput struct {
		Status string `json:"status"`
	}

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidSlotStatus(input.Status) {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	event.Status = models.SlotStatus(input.Status)
	event.UpdatedAt = time.Now()

	if err := database.DB.Save(&event).Error; err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func DeletEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 🔹 Extract event ID from URL
	  path := r.URL.Path
    // Assuming path format is /events/{id}
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        http.Error(w, "Invalid URL path", http.StatusBadRequest)
        return
    }
    idStr := parts[len(parts)-1]
    
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

	var event models.Event
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&event).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	
	if err := database.DB.Delete(&event).Error; err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Event deleted successfully")
}