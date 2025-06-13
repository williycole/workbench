package api

import (
	"log/slog"
	"net/http"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	slog.Info("󰊕 middlewareMetricsInc called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		hit := cfg.FileserverHits.Load()
		slog.Info("󰊕 middlewareMetricsInc hit", "hit", hit)
		next.ServeHTTP(w, r)
	})
}
