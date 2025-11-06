package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/models"
)

// SignupRequest is the JSON shape expected for signup.
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// simpleEmailRegexp is intentionally forgiving — it's fine for demo.
var simpleEmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateSignup validates the input and returns an error message if invalid.
func validateSignup(req SignupRequest) error {
	
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !simpleEmailRegexp.MatchString(req.Email) {
		return errors.New("invalid email")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// SignupHandler handles POST /signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := validateSignup(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if email already exists
	var existing models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&existing).Error; err == nil {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "server error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:         payload.Name,
		Email:        payload.Email,
		Password: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := SignupResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
