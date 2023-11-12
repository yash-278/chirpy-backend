package main

import (
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
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	fileServerhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileServerhandler)
	r.Handle("/app/*", fileServerhandler)

	adminRouter.Get("/metrics", apiCfg.handleMetrics)

	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/reset", apiCfg.reset)
	apiRouter.Post("/validate_chirp", validateChirp)

	r.Mount("/admin", adminRouter)
	r.Mount("/api", apiRouter)
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
