package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tiemouie01/chirpy/internal/auth"
	"github.com/tiemouie01/chirpy/internal/database"
)

type Chirp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func cleanChirp(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		for _, profaneWord := range profaneWords {
			if lowerWord == profaneWord {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	const characterError = "Chirp is too long"

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding JSON"+err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, characterError)
		return
	}

	// Get user ID from JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "You are not authorized to access this resource "+err.Error())
		return
	}

	// Clean chirp body
	cleanedBody := cleanChirp(params.Body)

	createChirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp")
		return
	}

	formattedChirp := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.UUID.String(),
	}

	respondWithJSON(w, 200, formattedChirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// Check for author_id in query params
	authorID := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if authorID != "" {
		// Parse UUID from string
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, 400, "Invalid authorID format")
			return
		}
		chirps, err = cfg.dbQueries.GetAllChirpsByAuthor(r.Context(), uuid.NullUUID{
			UUID:  authorUUID,
			Valid: true,
		})
		if err != nil {
			respondWithError(w, 500, "Error collecting author chirps")
		}
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, 500, "Error collecting chirps")
		return
	}

	formattedChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		formattedChirps[i] = Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.UUID.String(),
		}
	}

	respondWithJSON(w, 200, formattedChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Failed to parse ID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found.")
		return
	}

	formattedChirp := Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.UUID.String(),
	}

	respondWithJSON(w, 200, formattedChirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Ensure the user authorized to delete the chirp
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "You are not authorized to delete this chirp.")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 403, "You are not authorized to delete this chirp.")
		return
	}

	// Get the chirp ID from the token
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}

	// Delete the chirp from the database
	err = cfg.dbQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: id,
		UserID: uuid.NullUUID{
			UUID:  userId,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	w.WriteHeader(204)
}
