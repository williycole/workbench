package api

// TODO:
// 1. figure out a go api layout and how owyou would use your cfg correctly based on that layout
// 2. ask about the JWT handling, really understand that, the basics, what you typically do in most apps, etc
// 3. ask about how you can do routs like below
/*
TODO: and how this compares to how I was doing it
like how does this even work witout passing params in? so confusing...
mux.HandleFunc("/api/login", config.LoginUser)
mux.HandleFunc("/api/chirps", config.HandleChirps)
-- VS
mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})
*/
//TODO: ask about any other common patterns I may be missing

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

// UpgradeUserToChirpyRed is an endpoint that allows users to upgrade their account to Chirpy Red.
func (cfg *ApiConfig) UpgradeUserToChirpyRed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.PolkaKey || err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error: invalid or missing API key"})
		return
	}
	type upgradeUserToChirpyRedRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	var req upgradeUserToChirpyRedRequest
	d := json.NewDecoder(r.Body)
	err = d.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if req.Data.UserID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "missing user_id"})
		return
	}
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid user_id"})
		return
	}
	_, err = cfg.DbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}
	_, err = cfg.DbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to upgrade user"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// LoginUser handles user login, authenticating the user and returning a JWT token.
func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
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
	json.NewEncoder(w).Encode(UserResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (cfg *ApiConfig) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.createUser(w, r)
	case http.MethodPut:
		cfg.handleUsersUpdate(w, r)
	default:
		http.NotFound(w, r)
	}
}

// CreateUser handles the creation of a new user.
func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(UserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *ApiConfig) handleUsersUpdate(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(UserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *ApiConfig) HandleChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.createChirp(w, r)
	case http.MethodGet:
		cfg.getChirps(w, r)
	default:
		http.NotFound(w, r)
	}
}

// HandleChirpWithOptions handles requests for a single chirp, allowing retrieval and deletion.
func (cfg *ApiConfig) HandleChirpWithOptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg.getChirp(w, r)
	case http.MethodDelete:
		cfg.deleteChirp(w, r)
	default:
		http.NotFound(w, r)
	}
}

// getChirp retrieves a single chirp by its ID from the database.
func (cfg *ApiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching chirps"})
		return
	}
	slog.Info("🐦 get_chirp hit", "chirp", chirp.Body, "created_at", chirp.CreatedAt, "updated_at", chirp.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)
}

// deleteChirp handles the deletion of a chirp by its ID.
func (cfg *ApiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
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
	// get chirp id from URL path
	chirpParam := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	chirpID, err := uuid.Parse(chirpParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error parsing chirp ID"})
		return
	}
	// get chirp to check ownership
	chirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "can't find chirp"})
		return
	}
	// check if chirp belongs to user
	if chirp.UserID != userID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error: chirp does not belong to user"})
		return
	}
	err = cfg.DbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "delete chirp failed"})
		return
	}
	// Ok
	w.WriteHeader(http.StatusNoContent)
}

// createChirp handles the creation of a new chirp.
func (cfg *ApiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
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

// getChirps retrieves all chirps from the database.
func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	params := r.URL.Query()
	authorID := params.Get("author_id")
	// if authorID is provided, get chirps by author
	if authorID != "" {
		err := cfg.getChirpsByAuthor(w, r, authorID, &chirps)
		if err != nil {
			return
		}
	} else {
		err := cfg.getAllTheChirps(w, r, &chirps)
		if err != nil {
			return
		}
	}
	// sort chirps by created_at
	sortOrder := params.Get("sort")
	switch sortOrder {
	case "", "asc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
		})
	}
	slog.Info("🐦🐦🐦 get_all_chirps hit", "count", len(chirps))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}

// getAllTheChirps retrieves all chirps from the database.
func (cfg *ApiConfig) getAllTheChirps(w http.ResponseWriter, r *http.Request, chirps *[]database.Chirp) error {
	// get all chirps
	c, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching chirps"})
		return err
	}
	*chirps = c
	return nil
}

// getChirpsByAuthor retrieves all chirps by a specific author from the database.
func (cfg *ApiConfig) getChirpsByAuthor(w http.ResponseWriter, r *http.Request, authorID string,
	chirps *[]database.Chirp,
) error {
	id, err := uuid.Parse(authorID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid user_id"})
	}
	c, err := cfg.DbQueries.GetAllChirpsByAuthor(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error fetching chirps for user"})
	}
	*chirps = c
	return nil
}

// RefreshToken handles the refresh of a JWT token using a refresh token.
func (cfg *ApiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
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
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
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

// ResetHits resets the file server hits counter.
func (cfg *ApiConfig) ResetHits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
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
