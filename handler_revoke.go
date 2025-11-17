package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/praneeth-ayla/Chirpy/internal/auth"
	"github.com/praneeth-ayla/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized, please login again", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Token:     refreshToken,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized, please login again", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
