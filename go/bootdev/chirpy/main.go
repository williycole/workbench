package main

import (
	"chirpy/api"
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	filepathRoot = http.Dir(".")
	port         = "8080"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("⚠️ Error connecting to database: %v", err))
	}
	dbQueries := database.New(db)
	jwt := os.Getenv("JWT_SECRET")
	config := api.ApiConfig{DbQueries: dbQueries, Platform: platform, JWTSecret: jwt}
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
	// -- API Routes
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		config.LoginUser(w, r)
	})
	mux.HandleFunc("/api/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		config.RefreshToken(w, r)
	})
	mux.HandleFunc("/api/revoke", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		config.RevokeRefreshToken(w, r)
	})

	mux.HandleFunc("/api/chirps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			config.CreateChirp(w, r)
		case http.MethodGet:
			config.GetAllChirps(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/chirps/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			config.GetChirp(w, r)
			return
		}
		http.NotFound(w, r)
	})
	// -- Admin Routes
	mux.Handle("GET /admin/metrics", http.HandlerFunc(config.FileServerHitsHandler))
	mux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		config.ResetHits(w, r)
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			config.CreateUser(w, r)
		case http.MethodPut:
			config.HandlerUsersUpdate(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	// -- App Routes
	mux.Handle("/app/", config.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(filepathRoot))))
	mux.Handle("/app/assets/", config.MiddlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
