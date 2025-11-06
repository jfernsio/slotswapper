package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/models"
)

// ✅ 1️⃣ GET /api/swappable-slots
func GetSwappableSlots(w http.ResponseWriter, r *http.Request) {
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
	if err := database.DB.
		Where("user_id != ? AND status = ?", userID, models.SlotSwappable).
		Find(&events).Error; err != nil {
		http.Error(w, "Error fetching swappable slots", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// ✅ 2️⃣ POST /api/swap-request
func CreateSwapRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type RequestInput struct {
		MySlotID    uint `json:"mySlotId"`
		TheirSlotID uint `json:"theirSlotId"`
	}

	var input RequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var mySlot, theirSlot models.Event
	if err := database.DB.First(&mySlot, input.MySlotID).Error; err != nil {
		http.Error(w, "My slot not found", http.StatusNotFound)
		return
	}
	if err := database.DB.First(&theirSlot, input.TheirSlotID).Error; err != nil {
		http.Error(w, "Their slot not found", http.StatusNotFound)
		return
	}

	// Verify both are swappable
	if mySlot.Status != models.SlotSwappable || theirSlot.Status != models.SlotSwappable {
		http.Error(w, "Both slots must be swappable", http.StatusBadRequest)
		return
	}

	swap := models.SwapRequest{
		MySlotID:    input.MySlotID,
		TheirSlotID: input.TheirSlotID,
		RequesterID: mySlot.UserID,
		ReceiverID:  theirSlot.UserID,
		Status:      models.SwapPending,
	}

	if err := database.DB.Create(&swap).Error; err != nil {
		http.Error(w, "Failed to create swap request", http.StatusInternalServerError)
		return
	}

	// Lock both slots
	mySlot.Status = models.SlotSwapPending
	theirSlot.Status = models.SlotSwapPending
	database.DB.Save(&mySlot)
	database.DB.Save(&theirSlot)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(swap)
}

// ✅ 3️⃣ POST /api/swap-response?id=<swapID>
func RespondToSwap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ?id=<swapID> from query
	swapID := strings.TrimSpace(r.URL.Query().Get("id"))
	if swapID == "" {
		http.Error(w, "Missing swap ID", http.StatusBadRequest)
		return
	}

	type ResponseInput struct {
		Accept bool `json:"accept"`
	}

	var input ResponseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var swap models.SwapRequest
	if err := database.DB.First(&swap, swapID).Error; err != nil {
		http.Error(w, "Swap request not found", http.StatusNotFound)
		return
	}

	var mySlot, theirSlot models.Event
	database.DB.First(&mySlot, swap.MySlotID)
	database.DB.First(&theirSlot, swap.TheirSlotID)

	if input.Accept {
		swap.Status = models.SwapAccepted

		// Swap ownership
		temp := mySlot.UserID
		mySlot.UserID = theirSlot.UserID
		theirSlot.UserID = temp

		mySlot.Status = models.SlotBusy
		theirSlot.Status = models.SlotBusy
	} else {
		swap.Status = models.SwapRejected
		mySlot.Status = models.SlotSwappable
		theirSlot.Status = models.SlotSwappable
	}

	database.DB.Save(&swap)
	database.DB.Save(&mySlot)
	database.DB.Save(&theirSlot)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(swap)
}
