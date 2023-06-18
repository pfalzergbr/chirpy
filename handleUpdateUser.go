package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pfalzergbr/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), bcrypt.DefaultCost)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save user data")
		return
	}



	userUpdate := database.User{
		Id:       tokenClaims.Id,
		Email:    userParams.Email,
		Password: string(hashedPassword),
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

	respondWithJSON(w, http.StatusOK, sanitizedUser)
}
