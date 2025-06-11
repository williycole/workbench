package api

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"text/template"
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
}

// FileServerHitsHandler serves the file server hits page.
func (cfg *ApiConfig) FileServerHitsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	v := cfg.FileserverHits.Load()
	tmpl := `
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited {{.}} times!</p>
			</body>
		</html>
		`
	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		panic(fmt.Sprintf("⚠️ Error parsing template: %v", err))
	}
	err = t.Execute(w, v)
	if err != nil {
		slog.Error("Error executing template", "error", err)
	}
}

// ResetHits resets the file server hits counter.
func (cfg *ApiConfig) ResetHits(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "forbidden"})
		return
	}
	cfg.FileserverHits.Store(0)
	err := cfg.DbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		slog.Error("DeleteAllUsers failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error deleting users"})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset to 0"))
}

// CreateUser handles the creation of a new user.
func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req createUserRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}
	// hash password
	hp, err := auth.HashPassword(req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error hashing password"})
		return
	}
	// create user
	user, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: req.Email, HashedPassword: hp})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating user"})
		return
	}
	slog.Info("🧑 create_user hit", "email", user.Email, "created_at", user.CreatedAt, "updated_at", user.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// encode the user but ⚠️ WITHOUT the password
	json.NewEncoder(w).Encode(struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

// CreateChirp handles the creation of a new chirp.
func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body string `json:"body"`
	}
	var req createChirpRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}
	// Authenticate user via JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing or malformed Authorization header"})
		return
	}
	// Validate JWT
	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid or missing JWT"})
		return
	}
	// validate chirp/steralzie chirp
	cleanedBody, err := cfg.SteralizeChirp(req.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
	// create chirp
	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: userID})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating chirp"})
		return
	}
	slog.Info("🐦 create_chirp hit", "chirp", chirp.Body, "created_at", chirp.CreatedAt, "updated_at", chirp.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirp)
}

// GetAllChirps retrieves all chirps from the database.
func (cfg *ApiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching chirps"})
		return
	}
	slog.Info("🐦🐦🐦 get_all_chirps hit", "count", len(chirps))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}

// GetChirp retrieves a single chirp by its ID from the database.
func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpParam := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	chirpID, err := uuid.Parse(chirpParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error parsing chirp ID"})
		return
	}
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching chirps"})
		return
	}
	slog.Info("🐦 get_chirp hit", "chirp", chirp.Body, "created_at", chirp.CreatedAt, "updated_at", chirp.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
}

// LoginUser handles user login, authenticating the user and returning a JWT token.
func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	var req loginRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}
	// get user
	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching user"})
		return
	}
	// check password
	err = auth.CheckPasswordHash(user.HashedPassword, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error checking password"})
		return
	}
	// make refesh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error creating refresh token"})
		return
	}
	// store refresh token in DB
	expiresAt := time.Now().Add(60 * 24 * time.Hour) // 60 days

	_, err = cfg.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error storing refresh token"})
		return
	}
	// make token
	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error creating token"})
		return
	}
	slog.Info("🧑 login_user hit", "email", user.Email, "created_at", user.CreatedAt, "updated_at", user.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token,omitempty"`
	}{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// RefreshToken handles the refresh of a JWT token using a refresh token.
func (cfg *ApiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing or malformed Authorization header"})
		return
	}

	rt, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil || rt.ExpiresAt.Before(time.Now()) || (rt.RevokedAt.Valid && !rt.RevokedAt.Time.IsZero()) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid or expired refresh token"})
		return
	}

	user, err := cfg.DbQueries.GetUserByID(r.Context(), rt.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to create access token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{Token: newAccessToken})
}

// RevokeRefreshToken revokes a refresh token, making it invalid for future use.
func (cfg *ApiConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing or malformed Authorization header"})
		return
	}

	rt, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil || (rt.RevokedAt.Valid && !rt.RevokedAt.Time.IsZero()) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid refresh token"})
		return
	}

	err = cfg.DbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to revoke refresh token"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *ApiConfig) HandlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing or malformed Authorization header"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid or missing JWT"})
		return
	}

	type updateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error hashing password"})
		return
	}

	user, err := cfg.DbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error updating user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
