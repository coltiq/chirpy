package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       dbChirp.ID,
		AuthorID: dbChirp.AuthorID,
		Body:     dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortType := r.URL.Query().Get("sort")
	chirps := []Chirp{}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirps")
		return
	}

	if authorID == "" {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				AuthorID: dbChirp.AuthorID,
				Body:     dbChirp.Body,
			})
		}
	} else {
		authorIDint, err := strconv.Atoi(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Invalid author ID")
			return
		}
		for _, dbChirp := range dbChirps {
			if dbChirp.AuthorID == authorIDint {
				chirps = append(chirps, Chirp{
					ID:       dbChirp.ID,
					AuthorID: dbChirp.AuthorID,
					Body:     dbChirp.Body,
				})
			}
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortType == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
