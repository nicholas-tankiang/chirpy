package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type SuccessResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respBody := ErrorResponse{
			Error: "Chirp is too long",
		}
		writeJSONResponse(w, 400, respBody)
		return
	}

	respBody := SuccessResponse{
		CleanedBody: cleanChirp(params.Body),
	}
	writeJSONResponse(w, 200, respBody)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON %s", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
