package main

import (
	"strings"
)

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
