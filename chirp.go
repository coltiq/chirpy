package main

import (
	"encoding/json"
	"errors"
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

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := d.db.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp couldn't be created")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	return cleanChirp(body, badWords), nil
}

func cleanChirp(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
