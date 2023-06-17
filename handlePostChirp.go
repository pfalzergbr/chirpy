package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func (cfg apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleanedBody, err := handleGetValidatedChirp(params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	chirp, err := cfg.db.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func cleanProfanities(body string, badWords map[string]struct{}) string {
	bodyWords := strings.Split(body, " ")
	for wIdx, w := range bodyWords {
		if _, ok := badWords[strings.ToLower(w)]; ok {
			bodyWords[wIdx] = "****"
		}
	}

	return strings.Join(bodyWords, " ")
}

func handleGetValidatedChirp(params parameters) (string, error) {
	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		return "", fmt.Errorf("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedBody := cleanProfanities(params.Body, badWords)

	return cleanedBody, nil
}
