package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/praneeth-ayla/Chirpy/internal/auth"
)

func (cfg *apiConfig) webhookChirpRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	err = dec.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "No Content", nil)
		return
	}

	_, err = cfg.db.UpdateChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
