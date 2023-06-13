package main

import (
	"log"
	"net/http"
)

type ApiConfig struct {
	fileServerHits int
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileServerHits++
	return next
}

func main() {
	const filepathRoot = "./"
	const port = "8080"

	cfg := &ApiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir("."))))))
	mux.HandleFunc("/healthz", handlerReadiness)

	corsMux := middleware(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Printf("Serving files from %s\n", filepathRoot)

	log.Fatal(srv.ListenAndServe())

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
