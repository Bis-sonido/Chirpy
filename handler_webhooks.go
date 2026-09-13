package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bis-sonido/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {

	authAPIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid APIKey")
		return
	}

	if authAPIKey == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid APIKey")
		return
	}

	type upgradeChirpyRedUser struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	var params upgradeChirpyRedUser
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User doesn't exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
