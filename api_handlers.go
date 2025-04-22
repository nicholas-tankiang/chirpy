package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
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

func errorResponse(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	jsonResponse(w, code, errorResponse{
		Error: msg,
	})
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type SuccessResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, 500, "Error decoding parameters", err)
		return
	}

	if len(params.Body) > 140 {
		errorResponse(w, 400, "Chirp is too long", nil)
		return
	}

	respBody := SuccessResponse{
		CleanedBody: cleanChirp(params.Body),
	}
	jsonResponse(w, 200, respBody)
}
