package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Bis-sonido/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/Bis-sonido/Chirpy/internal/auth"
)

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string    `json:"body"`
	}

	type createChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var params createChirpRequest
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Post must be less than 140 characters")
		return
	}

	splitBody := strings.Split(params.Body, " ")
	for i, word := range splitBody {
		for _, badWord := range badWords {
			if strings.ToLower(word) == strings.ToLower(badWord) {
				// Replace the bad word with asterisks
				splitBody[i] = "****"
			}
		}
	}
	params.Body = strings.Join(splitBody, " ")

	authBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}
	if authBearer == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	userID, err := auth.ValidateJWT(authBearer, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, createChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
