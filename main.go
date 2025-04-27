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
	apiCfg := apiConfig{}

	fs := http.FileServer(http.Dir("."))
	appFS := http.StripPrefix("/app/", fs)
	if err := godotenv.Load(); err != nil {
		log.Println("warning, no .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appFS))
	mux.Handle("/", apiCfg.middlewareMetricsInc(appFS))
	mux.Handle("/assets/", apiCfg.middlewareMetricsInc(appFS))
	mux.HandleFunc("/metrics", apiCfg.getHitsCounter)
	mux.HandleFunc("/reset", apiCfg.resetHitsCounter)
	mux.HandleFunc("/healthz", app.ServerReadiness)

	addr := ":" + port

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server faileD: %v", err)
	}

}
