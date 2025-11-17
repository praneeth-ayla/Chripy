package main

import (
	"encoding/json"
	"net/http"

	"github.com/praneeth-ayla/Chirpy/internal/auth"
	"github.com/praneeth-ayla/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	err := dec.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get the headers", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
