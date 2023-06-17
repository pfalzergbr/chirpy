package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirps, err := cfg.db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp")
		return
	}

	chirp, ok := chirps.Chirps[id]

	if !ok {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
