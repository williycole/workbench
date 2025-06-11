package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt and returns the hashed password as a string.
func HashPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	s := string(encryptedPassword)
	return s, nil
}

// CheckPasswordHash checks if the provided password matches the hashed password.
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

type CustomClaims struct {
	Claim string `json:"claim"`
	jwt.RegisteredClaims
}

// MakeJWT creates a JWT token for the given user ID and token secret.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create the Claims
	rc := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}
	// signing method, SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	st, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return st, nil
}

// ValidateJWT validates the given JWT token and returns the user ID.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// parse the token with the custom claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	// check supbject and issuer
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy" {
		return uuid.Nil, err
	}
	// parse the subject as a UUID
	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}
	// return the user ID
	return id, nil
}

// GetBearerToken returns the value of the Authorization header
func GetBearerToken(headers http.Header) (string, error) {
	t := headers.Get("Authorization")
	if t == "" {
		return "", errors.New("no Authorization header found")
	}
	token := strings.TrimPrefix(t, "Bearer ")
	return token, nil
}

// MakeBearerToken creates a Bearer token string from the given JWT token.
func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(key)
	return s, nil
}
