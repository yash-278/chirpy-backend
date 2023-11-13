package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yash-278/chirpy-backend/database"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	DB, _ := database.NewDB("database.json")
	const port = "8000"

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	// Wrap the middlewareDb function to match the expected signature
	dbMiddleware := func(next http.Handler) http.Handler {
		return middlewareDb(next, DB)
	}

	apiRouter.Use(dbMiddleware)

	fileServerhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileServerhandler)
	r.Handle("/app/*", fileServerhandler)

	adminRouter.Get("/metrics", apiCfg.handleMetrics)

	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/reset", apiCfg.reset)

	apiRouter.Post("/chirps", addChirp)
	apiRouter.Get("/chirps", getChirps)

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
