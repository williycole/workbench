package api

import (
	"chirpy/internal/database"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// RespBody is a simple struct to hold the response body.
type RespBody struct {
	Body string `json:"body"`
}

// ErrorResponse is a simple struct to hold error messages.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ApiConfig holds the configuration for the API, including the file server hits counter
type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	JWTSecret      string
	PolkaKey       string
}

// UserResponse is a struct that represents a user response.
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}
