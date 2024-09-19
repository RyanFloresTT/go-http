package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerGetMetricHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, cfg.fileServerHits)
	w.Write([]byte(html))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	input := chirpRequest{}

	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	if len(input.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "chirp is too long")
		return
	}

	checkProfanity(&input)

	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": input.Body})
}

type chirpRequest struct {
	Body string `json:"body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	data, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		fmt.Errorf(err.Error())
	}
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	w.Write(data)
}

func checkProfanity(input *chirpRequest) {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range profanity {
		loweredBody := strings.ToLower(input.Body)
		loweredWord := strings.ToLower(word)

		if strings.Contains(loweredBody, loweredWord) {
			startIndex := strings.Index(loweredBody, loweredWord)
			endIndex := startIndex + len(word)

			input.Body = input.Body[:startIndex] + "****" + input.Body[endIndex:]
		}
	}
}
