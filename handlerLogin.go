package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	users, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}

	for _, user := range dbUsers {
		if user.Email == params.Email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
				return
			}
			respondWithJSON(w, http.StatusOK, User{
				Id:    user.Id,
				Email: user.Email,
			})
			return
		}
	}

	respondWithError(w, http.StatusInternalServerError, "Not able to find Email")
}
