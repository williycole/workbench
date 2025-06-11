package api

import (
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"text/template"

	"github.com/google/uuid"
)

type RespBody struct {
	Body string `json:"body"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
}

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

// func (cfg *ApiConfig) ValidateChirp(w http.ResponseWriter, r *http.Request) {
// 	// clean up chirp
// 	req, err := cfg.SteralizeChirp(r)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
// 		return
// 	}
// 	// handle ok
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]any{
// 		"valid":        true,
// 		"cleaned_body": req.Body,
// 	})
//
// 	slog.Info("🐱 validate_chirp hit", "body", req.Body)
// }

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email string `json:"email"`
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
	user, err := cfg.DbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating user"})
		return
	}

	slog.Info("🧑 create_user hit", "email", user.Email, "created_at", user.CreatedAt, "updated_at", user.UpdatedAt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
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
	// parse userID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error parsing user ID"})
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
