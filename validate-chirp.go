package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type jsonError struct {
		Error string `json:"error"`
	}

	type jsonSuccess struct {
		Valid bool `json:"valid"`
	}

	errorRespBody := jsonError{}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		errorRespBody.Error = "Something went wrong"
		dat, err := json.Marshal(errorRespBody)
		if err != nil {
			errorRespBody.Error = "Something went wrong"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write(dat)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		errorRespBody.Error = "Chirp is too long"
		dat, err := json.Marshal(errorRespBody)
		if err != nil {
			errorRespBody.Error = "Something went wrong"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write(dat)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	validChirp := jsonSuccess{Valid: true}

	dat, err := json.Marshal(validChirp)
	if err != nil {
		errorRespBody.Error = "Something went wrong"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
