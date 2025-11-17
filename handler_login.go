package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/praneeth-ayla/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password          string `json:"password"`
		Email             string `json:"email"`
		ExpiriesInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expiry time.Duration
	if params.ExpiriesInSeconds > 0 {
		expiry = time.Duration(params.ExpiriesInSeconds) * time.Second

		if expiry > time.Hour {
			expiry = time.Hour
		}
	} else {
		expiry = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		},
	})
}
