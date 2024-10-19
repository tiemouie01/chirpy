package main

import "net/http"

type jsonResponse struct {
	Hits int    `json:"hits"`
	Msg  string `json:"msg"`
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Cannot reset application in production environment")
		return
	}

	// Reset file server hits to 0
	cfg.fileserverHits.Store(0)

	// Delete all users from the database
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error deleting users from database")
		return
	}

	response := jsonResponse{Hits: 0, Msg: "All users deleted from the database"}

	respondWithJSON(w, 200, response)
}
