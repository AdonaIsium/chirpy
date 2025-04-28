package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AdonaIsium/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
type ErrorResp struct {
	Error string `json:"error"`
}
type CleanedChirp struct {
	CleanedBody string `json:"cleaned_body"`
}
type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {
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

	cleanedChirp := CleanedChirp{
		CleanedBody: cleanedBody,
	}

	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedChirp.CleanedBody, UserID: chirp.UserID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating DB entry for chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})
}

func (cfg apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Getting all Chirps", err)
		return
	}
	resp := make([]ChirpResponse, len(chirps))
	for i, c := range chirps {
		resp[i] = ChirpResponse{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	if chirpIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing chirp ID", nil)
		return
	}
	fmt.Printf(chirpIDStr)

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to resolve ID into UUID", err)
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Getting Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
