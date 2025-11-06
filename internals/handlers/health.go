package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jfernsio/slotswapper/internals/database"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := database.DB.DB()
	if err != nil {
		http.Error(w, "db not ready", http.StatusInternalServerError)
		return
	}
	if err := sqlDB.Ping(); err != nil {
		http.Error(w, "db ping failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
