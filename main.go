package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AdonaIsium/chirpy/app"
	"github.com/joho/godotenv"
)

func main() {
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
	mux.Handle("/app/", appFS)
	mux.Handle("/", appFS)
	mux.Handle("/assets", appFS)
	mux.HandleFunc("/healthz", app.ServerReadiness)

	addr := ":" + port

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server faileD: %v", err)
	}

}
