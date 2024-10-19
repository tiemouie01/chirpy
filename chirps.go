package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	const defaultError = "Something went wrong"
	const characterError = "Chirp is too long"

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, defaultError)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, characterError)
		return
	}

	cleanedBody := cleanChirp(params.Body)
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID")
		return
	}

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
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error collecting chirps.")
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
