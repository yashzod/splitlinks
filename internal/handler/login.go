package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
}

func SignUpHandler(db *gorm.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Validate minimal required fields
		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}

		// Check if email already exists
		var existing User
		if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}

		// Hash password
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "could not hash password", http.StatusInternalServerError)
			return
		}

		// Create user
		user := User{
			ID:           uuid.New(),
			Email:        req.Email,
			FullName:     req.FullName,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			ImageURL:     req.ImageURL,
			PasswordHash: string(hashedPwd),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "could not create user", http.StatusInternalServerError)
			return
		}

		// Generate JWT
		tokenStr, err := MakeJWT(user.ID.String(), jwtSecret, 24*time.Hour)
		if err != nil {
			http.Error(w, "could not create token", http.StatusInternalServerError)
			return
		}

		// Return response
		var resp LoginResponse
		resp.Token = tokenStr
		resp.User.ID = user.ID.String()
		resp.User.Email = user.Email
		resp.User.FullName = user.FullName
		resp.User.ImageURL = user.ImageURL
		resp.User.CreatedAt = user.CreatedAt

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
