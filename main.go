package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/pfalzergbr/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
}

func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB(filepathRoot + "/database.json")
	if err != nil {
		log.Fatalf("Error creating database: %s\n", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	// * API Routes
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/chirps", apiCfg.handleGetChirps)
	apiRouter.Get("/chirps/{id}", apiCfg.handleGetChirp)
	apiRouter.Post("/chirps", apiCfg.handlePostChirp)

	apiRouter.Post("/users", apiCfg.handleCreateUser)
	apiRouter.Post("/login", apiCfg.handleLoginUser)

	// * Admin Routes
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
