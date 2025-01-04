package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/tiemouie01/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  Data
	}

	// Ensure the requesting resource is authenticated
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if apiKey != cfg.polkaApiKey {
		respondWithError(w, 401, "You are not authorized to access this API")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding JSON")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// Mark the user as a Chirpy Red member
	err = cfg.dbQueries.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	w.WriteHeader(204)
}
