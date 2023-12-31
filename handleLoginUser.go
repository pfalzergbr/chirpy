package main

import (
	"encoding/json"
	"net/http"

	"github.com/pfalzergbr/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userParams := userBody{}
	err := decoder.Decode(&userParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(userParams.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userParams.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}


	jwt, err := cfg.createJWT(user.Id, userParams.ExipiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not log in")
		return
	}

	userResponse := database.LoginResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: jwt,
	}

	respondWithJSON(w, http.StatusOK, userResponse)

}
