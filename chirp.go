package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (d *dbConfig) ChirpGetHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := d.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (d *dbConfig) ChirpPostHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanChirp(params.Body)
	chirp, err := d.db.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp couldn't be created")
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func cleanChirp(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")

	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
