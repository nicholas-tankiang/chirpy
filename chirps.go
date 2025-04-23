package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, 500, "Error decoding parameters", err)
		return
	}

	// validate chirp
	if len(params.Body) > 140 {
		errorResponse(w, 400, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanChirp(params.Body),
			UserID: params.UserID,
		},
	)
	if err != nil {
		errorResponse(w, 500, "Error creating chirp", err)
		return
	}

	jsonResponse(w, 201, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	rawChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error retrieving chirps", err)
		return
	}

	//map database model format []database.Chirp to the API model format of JSON field names
	responseChirps := []Chirp{}
	for _, dbChirp := range rawChirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	jsonResponse(w, 200, responseChirps)
}

// converts words in map to censor
func cleanChirp(body string) string {
	censor := "****"
	profaneWordsMap := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	splitBody := strings.Split(body, " ")
	for i := range splitBody {
		if profaneWordsMap[strings.ToLower(splitBody[i])] {
			splitBody[i] = censor
		}
	}
	cleanedBody := strings.Join(splitBody, " ")

	return cleanedBody
}
