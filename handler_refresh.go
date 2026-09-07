package main

import (
	"net/http"
	"time"

	"github.com/Bis-sonido/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	authBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}
	if authBearer == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), authBearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}

	type refreshTokenResponse struct {
		Token string `json:"token"`
	}
	
	respondWithJSON(w, http.StatusOK, refreshTokenResponse{
		Token: token,
	})
}