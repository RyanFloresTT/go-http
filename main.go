package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	port := "8080"

	var apiCfg apiConfig
	apiCfg.fileServerHits = 0

	setupServeMux(mux, &apiCfg)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s\n", server.Addr)

	server.ListenAndServe()
}

type apiConfig struct {
	fileServerHits int
}

func setupServeMux(mux *http.ServeMux, apiCfg *apiConfig) {
	fileServer := http.FileServer(http.Dir("./app"))
	stripPrefixHandler := http.StripPrefix("/app", fileServer)
	wrappedHandler := apiCfg.middlewareCountHits(stripPrefixHandler)

	mux.Handle("/app/", wrappedHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerGetMetricHits)
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("/api/reset", apiCfg.handlerResetHits)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
}
