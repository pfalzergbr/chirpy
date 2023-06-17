package main

import (
	"net/http"
	"sort"

	"github.com/pfalzergbr/chirpy/internal/database"
)

func (cfg apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbStruct, err := cfg.db.GetChirps()

	chirps := make([]database.Chirp, 0, len(dbStruct.Chirps))

	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, http.StatusOK, chirps)
}
