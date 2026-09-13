package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("Error hashing password")
		return "", err
	}

	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Println("Error verifying password")
		return false, err
	}

	return match, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	key := []byte(tokenSecret)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	key := []byte(tokenSecret)

	claimStruct := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claimStruct, func(token *jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	uuidSubject, err := uuid.Parse(claimStruct.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return uuidSubject, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	strippedToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if strippedToken == "" {
		return "", fmt.Errorf("missing Bearer token")
	}

	return strippedToken, nil
}

func GetAPIKey(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	strippedAPIKey := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))
	if strippedAPIKey == "" {
		return "", fmt.Errorf("missing APIKey")
	}
	return strippedAPIKey, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}
