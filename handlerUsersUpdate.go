package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/coltiq/chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
		return
	}
	userToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(userToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return cfg.jwtSecret, nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var userID int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["sub"].(string); ok {
			userID, err = strconv.Atoi(id)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Unable to get user id")
				return
			}
		} else {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}
	} else {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	user, err := cfg.DB.UpdateUser(userID, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    userID,
			Email: user.Email,
		},
	})
}
