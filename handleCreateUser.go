package main

import (
	"encoding/json"
	"net/http"
)

func (cfg apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type userBody struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userParams := userBody{}
	err := decoder.Decode(&userParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.CreateUser(userParams.Email, userParams.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
