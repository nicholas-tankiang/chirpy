package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// User struct embedded in a response in case of future additions
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, 500, "Error decoding parameters", err)
		return
	}

	//cannot directly pass the hashpassword output to CreateUser, storing in hash
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to hash password", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	},
	)
	if err != nil {
		errorResponse(w, 500, "Error creating user", err)
		return
	}

	jsonResponse(w, 201, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		errorResponse(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}

	//reset metrics from metrics.go
	cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "All users deleted successfully, fileserverHits reset to zero",
	})
}
