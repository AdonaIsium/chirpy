package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
}
type ErrorResp struct {
	Error string `json:"error"`
}
type ValidResp struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(chirp.Body) > 140 {
		log.Printf("Chirp must be less than 140 chars. Received %d", len(chirp.Body))
		errorResp := ErrorResp{
			Error: "Chirp is too long",
		}
		jsonResp, err := json.Marshal(errorResp)
		if err != nil {
			log.Printf("Error encoding error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

	splitChirp := strings.Split(chirp.Body, " ")
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range splitChirp {
		for _, profanity := range profaneWords {
			if strings.ToLower(word) == profanity {
				splitChirp[i] = "****"
			}
		}
	}
	cleanedBody := strings.Join(splitChirp, " ")

	validResp := ValidResp{
		CleanedBody: cleanedBody,
	}

	jsonResp, err := json.Marshal(validResp)
	if err != nil {
		log.Printf("Error encoding valid: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
