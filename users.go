package main

import (
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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type paramters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)

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
