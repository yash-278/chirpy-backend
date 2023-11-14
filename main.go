package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yash-278/chirpy-backend/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	DB, _ := database.NewDB("database.json")

	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             DB,
	}

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	fileServerhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileServerhandler)
	r.Handle("/app/*", fileServerhandler)

	adminRouter.Get("/metrics", apiCfg.handleMetrics)

	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/reset", apiCfg.reset)

	apiRouter.Post("/chirps", apiCfg.addChirp)
	apiRouter.Get("/chirps", apiCfg.getChirps)

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
