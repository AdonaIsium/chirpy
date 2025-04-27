package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}
type ErrorResp struct {
	Error string `json:"error"`
}
type ValidResp struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirp.Body) > 140 {
		log.Printf("Chirp must be less than 140 chars. Receied %d", len(chirp.Body))
		errorResp := ErrorResp{
			Error: "Chirp is too long",
		}
		jsonResp, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error encoding error: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonResp)
		return
	}

	validResp := ValidResp{
		Valid: true,
	}
	jsonResp, err := json.Marshal(validResp)
	if err != nil {
		log.Printf("Error encoding valid: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
