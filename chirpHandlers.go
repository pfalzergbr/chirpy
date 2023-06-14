package main

import (
	"encoding/json"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpBody struct {
		Body string `json:"body"`
	}

	// type chirpValidationResp struct {
	// 	Error string `json:"error"`
	// 	Valid bool   `json:"valid"`
	// }

	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}

	err := decoder.Decode(&chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid": true}`))
}
