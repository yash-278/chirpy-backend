package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/yash-278/chirpy-backend/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	secretKey      string
	polkaKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	DB, _ := database.NewDB("database.json")

	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             DB,
		secretKey:      secretKey,
		polkaKey:       polkaKey,
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
	apiRouter.Get("/chirps/{id}", apiCfg.getChirp)
	apiRouter.Delete("/chirps/{id}", apiCfg.handlerDeleteChirp)

	apiRouter.Post("/users", apiCfg.addUser)
	apiRouter.Put("/users", apiCfg.handlerUserUpdate)

	apiRouter.Post("/refresh", apiCfg.handleRefreshTokenCreate)
	apiRouter.Post("/revoke", apiCfg.handleRefreshTokenRevoke)

	apiRouter.Post("/login", apiCfg.loginUser)

	apiRouter.Post("/polka/webhooks", apiCfg.handlerWebhooks)

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
