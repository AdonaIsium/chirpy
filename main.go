package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AdonaIsium/chirpy/app"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHitsCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	count := cfg.fileserverHits.Load()
	fmt.Fprintf(w, `Hits: %d`, count)
}

func (cfg *apiConfig) resetHitsCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte(`"Metrics reset"`))
}

func main() {
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	if err := godotenv.Load(); err != nil {
		log.Println("warning, no .env file found")
	}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHitsCounter)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)
	mux.HandleFunc("GET /api/healthz", app.ServerReadiness)

	addr := ":" + port

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server faileD: %v", err)
	}

}
