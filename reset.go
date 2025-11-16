package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.ResetDB(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to reset", err)
		return
	}

	cfg.fileserverHits.Store(0)

	respondWithJSON(w, 200, nil)
}
