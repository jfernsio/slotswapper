package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/models"
	"github.com/jfernsio/slotswapper/internals/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashed)

	if err := database.DB.Create(&input).Error; err != nil {
		http.Error(w, "User already exists or invalid data", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := database.DB.Where("email = ?", input.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		http.Error(w, "User does not exist!", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
