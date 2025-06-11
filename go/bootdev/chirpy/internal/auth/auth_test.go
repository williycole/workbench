package auth

import (
	"net/http"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Can create JWT and claims are correct",
			userID:      uuid.New(),
			tokenSecret: "foo-bar",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}

			// Parse the token and check claims
			token, err := jwt.ParseWithClaims(got, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
				return []byte(tt.tokenSecret), nil
			})
			if err != nil {
				t.Fatalf("Failed to parse generated JWT: %v", err)
			}
			claims, ok := token.Claims.(*jwt.RegisteredClaims)
			if !ok {
				t.Fatalf("Claims are not of type *jwt.RegisteredClaims")
			}
			if claims.Issuer != "chirpy" {
				t.Errorf("Issuer = %v, want %v", claims.Issuer, "chirpy")
			}
			if claims.Subject != tt.userID.String() {
				t.Errorf("Subject = %v, want %v", claims.Subject, tt.userID.String())
			}
			if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
				t.Errorf("ExpiresAt is invalid or in the past")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	// Generate a valid token for positive test case
	validToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create valid JWT for test: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token returns correct UUID",
			tokenString: validToken,
			tokenSecret: tokenSecret,
			want:        userID,
			wantErr:     false,
		},
		{
			name:        "Invalid secret returns error",
			tokenString: validToken,
			tokenSecret: "wrong-secret",
			want:        uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Malformed token returns error",
			tokenString: "not.a.jwt",
			tokenSecret: tokenSecret,
			want:        uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Bearer token",
			headers: http.Header{"Authorization": []string{"Bearer mytoken123"}},
			want:    "mytoken123",
			wantErr: false,
		},
		{
			name:    "No Authorization header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Authorization header without Bearer prefix",
			headers: http.Header{"Authorization": []string{"Basic something"}},
			want:    "Basic something",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
