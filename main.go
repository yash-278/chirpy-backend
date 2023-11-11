package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	const port = "8000"

	r := chi.NewRouter()
	mux := http.NewServeMux()

	fileServerhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileServerhandler)
	r.Handle("/app/*", fileServerhandler)

	r.Get("/healthz", handleReadiness)

	r.Get("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.reset)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metric := fmt.Sprintf("Hits: %v", cfg.fileserverHits)
	w.Write([]byte(metric))
}
