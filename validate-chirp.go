package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type jsonError struct {
	Error string `json:"error"`
}

type jsonSuccess struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	errorRespBody := jsonError{Error: msg}
	dat, err := json.Marshal(errorRespBody)
	if err != nil {
		errorRespBody.Error = "Something went wrong"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(dat)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload jsonSuccess) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	respondWithJSON(w, 200, jsonSuccess{CleanedBody: cleanedBody})
}
