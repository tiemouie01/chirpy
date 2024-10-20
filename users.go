package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tiemouie01/chirpy/internal/auth"
	"github.com/tiemouie01/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type paramters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing user password")
		return
	}

	// Create special params to satisfy db query
	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)

	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}

	formattedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, formattedUser)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type paramters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	// Search for user in db
	user, err := cfg.dbQueries.FindUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Failed to fetch user")
		return
	}

	// Check if db password matches password in params
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, 403, "Invalid login credentials")
		return
	}

	// If passwords match, return user details
	formattedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 200, formattedUser)
}
