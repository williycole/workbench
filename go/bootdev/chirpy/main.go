package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	const filepathRoot = http.Dir(".")

	// ServeMux in Go indeed acts as an orchestrator or router for incoming HTTP requests. It's responsible for directing each request to the appropriate handler
	routerMux := http.NewServeMux()

	// http.Server allows us to define ther server's characteristics
	// including hook in our server's handler, ie ServeMux or NewServeMux
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        routerMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	routerMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	routerMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(filepathRoot)))
	routerMux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets"))))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
