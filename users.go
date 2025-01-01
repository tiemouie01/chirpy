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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	// Create user JWT
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 3600)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	// Create Refresh Token and store it in the database.
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	dbToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Second * 60 * 60 * 7 * 30 * 2),
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// If passwords match, return user details
	formattedUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: dbToken.Token,
	}
	respondWithJSON(w, 200, formattedUser)
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get old token from header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	// Look up the token in the database
	userId, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	// Create the refresh token
	token, err := auth.MakeJWT(userId, cfg.jwtSecret, 3600)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// Return the token
	type Token struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, Token{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	// Get the token from the header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	// Revoke the token
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	// Return 204 status code
	w.WriteHeader(204)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the json parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "There was an error handling the request")
		return
	}

	// Get the user from using their access token.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "You are not authorized to perform this action.")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "You are not authorized to perform this action.")
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(params.Email)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// Update the user record in the database
	user, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:    userId,
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	})
	if err != nil {
		respondWithError(w, 500, "Error updating the user's information.")
		return
	}

	type UpdatedUser struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	respondWithJSON(w, 200, UpdatedUser{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
