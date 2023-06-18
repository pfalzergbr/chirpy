package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	Id        int       `json:"id"`
	ExpiresAt time.Time `json:"exp"`
}

func (cfg apiConfig) createJWT(id int, expiresAt *int) (string, error) {
	var expirationTime time.Time

	if expiresAt == nil {
		expirationTime = time.Now().Add(time.Duration(*expiresAt) * time.Second)
	} else if *expiresAt > 60*60*24 {
		expirationTime = time.Now().Add(24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(24 * time.Hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.Itoa(id),
	})

	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (cfg apiConfig) validateJWT(tokenString string) (tokenClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})

	if err != nil {
		return tokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.RegisteredClaims); ok && token.Valid {
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return tokenClaims{}, err
		}

		return tokenClaims{
			Id:        id,
			ExpiresAt: claims.ExpiresAt.Time,
		}, nil
	}

	return tokenClaims{}, nil
}
