package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"net/http"
)

func TestHashPassword(t *testing.T) {
	password := "isThisPssw0rdG00d?"

	// hashing succeeds and changed the plain text
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("expecting no error, got %v", err)
	}
	if hash == "" {
		t.Error("expected a non-empty hash string")
	}
	if hash == password {
		t.Error("hash should not match plain text")
	}

	// validation succeeds with the exact same password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("expecting no error during password check:  got %v", err)
	}
	if !match {
		t.Error("expected correct password to match but failed")
	}

	// validation fails with the wrong password
	wrongMatch, err := CheckPasswordHash("This1sW4r9ng", hash)
	if err != nil {
		t.Errorf("expected no error during a failed password check, got: %v", err)
	}
	if wrongMatch {
		t.Error("incorrect password successful")
	}
}

func TestCreateAndValidateJWT(t *testing.T) {

	key := "thisIsASecretKey"
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	// create a JWT
	token, err := MakeJWT(userID, key, expiresIn)
	if err != nil {
		t.Fatalf("expected no error during JWT creation, got: %v", err)
	}
	if token == "" {
		t.Error("expected a non-empty JWT string")
	}

	// validate the JWT
	validatedUserID, err := ValidateJWT(token, key)
	if err != nil {
		t.Fatalf("expected no error during JWT validation, got: %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("expected user ID %v, got %v", userID, validatedUserID)
	}

	// expired JWT should fail validation
	expiredToken, err := MakeJWT(userID, key, -1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error during JWT creation, got: %v", err)
	}
	_, err = ValidateJWT(expiredToken, key)
	if err == nil {
		t.Error("expected an error during expired JWT validation, got none")
	}

	// wrong secret key should fail validation
	_, err = ValidateJWT(token, "wrongSecretKey")
	if err == nil {
		t.Error("expected an error during JWT validation with wrong secret key, got none")
	}
}

func TestValidateGetBearerToken(t *testing.T) {
	
	type testCase struct {
		name string
		headers http.Header
		expectedToken string
		expectError bool
	}

	testCases := []testCase{
		{
			name: "Valid Bearer Token",
			headers: http.Header{"Authorization": []string{"Bearer validToken"}},
			expectedToken: "validToken",
			expectError: false,
		},
		{
			name: "Missing Authorization Header",
			headers: http.Header{},
			expectedToken: "",
			expectError: true,
		},
		{
			name: "Invalid Token Format",
			headers: http.Header{"Authorization": []string{"InvalidFormat"}},
			expectedToken: "",
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if err != nil {
				if !tt.expectError {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if tt.expectError {
					t.Error("expected an error but got none")
				}
				if token != tt.expectedToken {
					t.Errorf("expected token %s, got %s", tt.expectedToken, token)
				}
			}
		})
	}
}