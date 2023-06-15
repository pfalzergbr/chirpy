package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var profanities = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func handleChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	cleanedBody := cleanProfanities(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
}

func cleanProfanities(body string) string {
	bodyWords := strings.Split(body, " ")
	for wIdx, w := range bodyWords {
		for _, p := range profanities {
			if strings.EqualFold(w, p) {
				bodyWords[wIdx] = "****"
			}
		}
	}

	return strings.Join(bodyWords, " ")
}
