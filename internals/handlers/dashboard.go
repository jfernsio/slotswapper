package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/models"
)

// func Profile(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	userID := r.Context().Value("user_id")
// 	if userID == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var user models.User
// 	if err := database.DB.First(&user, userID).Error; err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}



// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get user info
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	//only send back user name and id
	userRes := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
	}

	// get events for this user
	var events []models.Event
	if err := database.DB.Where("user_id = ?", user.ID).Find(&events).Error; err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}

	// combine both in one response
	response := map[string]interface{}{
		"user":   userRes,
		"events": events,
	}

	json.NewEncoder(w).Encode(response)
}