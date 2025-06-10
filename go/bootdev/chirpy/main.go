package main

import (
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

type RespBody struct {
	Body string `json:"body"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	slog.Info("󰊕 middlewareMetricsInc called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		hit := cfg.fileserverHits.Load()
		slog.Info("󰊕 middlewareMetricsInc hit", "hit", hit)
		next.ServeHTTP(w, r)
	})
}

const (
	filepathRoot = http.Dir(".")
	port         = "8080"
)

func main() {
	os.Setenv("PLATFORM", os.Getenv("PLATFORM"))
	os.Setenv("DB_URL", os.Getenv("DATABASE_URL"))

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("⚠️ Error connecting to database: %v", err))
	}
	dbQueries := database.New(db)
	config := apiConfig{dbQueries: dbQueries, platform: os.Getenv("PLATFORM")}

	// ServeMux in Go indeed acts as an orchestrator or router for incoming HTTP requests. It's responsible for directing each request to the appropriate handler
	mux := http.NewServeMux()

	// http.Server allows us to define ther server's characteristics
	// including hook in our server's handler, ie ServeMux or NewServeMux
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		validateChirp(w, r)
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		config.createUser(w, r)
	})

	mux.Handle("GET /admin/metrics", http.HandlerFunc(config.fileServerHitsHandler))
	mux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		config.resetHitsHandler(w, r)
	})

	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(filepathRoot))))
	mux.Handle("/app/assets/", config.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}

func (cfg *apiConfig) fileServerHitsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	v := cfg.fileserverHits.Load()

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

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "forbidden"})
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		slog.Error("DeleteAllUsers failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "error deleting users"})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset to 0"))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	// check method TODO: this really should happen before this func gets called
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}
	// decode the request body
	var req RespBody
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error decoding request body"})
		return
	}
	// check length
	if len(req.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"})
		return
	}
	// handle dirty words
	dirtyWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range dirtyWords {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		req.Body = re.ReplaceAllString(req.Body, "****")
	}
	// handle ok
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"valid":        true,
		"cleaned_body": req.Body,
	})

	slog.Info("🐱 validate_chirp hit", "body", req.Body)
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
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
