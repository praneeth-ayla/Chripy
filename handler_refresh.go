package main

import (
	"net/http"
	"time"

	"github.com/praneeth-ayla/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized, please login again", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized, please login again", err)
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: jwtToken})
}
