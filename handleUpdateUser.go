package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pfalzergbr/chirpy/internal/database"
)

func (cfg apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userParams := userBody{}
	err := decoder.Decode(&userParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	auth := r.Header.Get("Authorization")
	token := strings.Split(auth, " ")[1]

	tokenClaims, err := cfg.validateJWT(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	now := time.Now()
	if tokenClaims.ExpiresAt.Unix() < now.Unix() {
		respondWithError(w, http.StatusUnauthorized, "Token expired")
		return
	}

	userUpdate := database.User{
		Id:       tokenClaims.Id,
		Email:    userParams.Email,
		Password: userParams.Password,
	}

	user, err := cfg.db.UpdateUser(tokenClaims.Id, userUpdate)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	sanitizedUser := database.CreateUserResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, sanitizedUser)
}
