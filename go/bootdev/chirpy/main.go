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

	polkaKey := os.Getenv("POLKA_KEY")
	cfg := api.ApiConfig{DbQueries: dbQueries, Platform: platform, JWTSecret: jwt, PolkaKey: polkaKey}
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
	// -- Api Routes
	mux.HandleFunc("/api/login", cfg.LoginUser)
	mux.HandleFunc("/api/refresh", cfg.RefreshToken)
	mux.HandleFunc("/api/revoke", cfg.RevokeRefreshToken)
	mux.HandleFunc("/api/chirps", cfg.HandleChirps)
	mux.HandleFunc("/api/chirps/", cfg.HandleChirpWithOptions)
	mux.HandleFunc("/api/users", cfg.HandleUsers)
	mux.HandleFunc("/api/polka/webhooks", cfg.UpgradeUserToChirpyRed)
	// -- Admin Routes
	mux.Handle("GET /admin/metrics", http.HandlerFunc(cfg.FileServerHitsHandler))
	mux.HandleFunc("/admin/reset", cfg.ResetHits)
	// -- App Routes
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(filepathRoot))))
	mux.Handle("/app/assets/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
