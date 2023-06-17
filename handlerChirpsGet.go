package main

import (
	"net/http"
	"sort"
)

func (cfg apiConfig) HandleGetChiprs(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}
	
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })


	respondWithJSON(w, http.StatusOK, chirps)
}
