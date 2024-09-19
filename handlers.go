package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		errorResponse := errorResponse{
			Error: "something went wrong",
		}

		data, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		w.WriteHeader(500)
		w.Write(data)
		return
	}

	if len(input.Body) > 140 {
		errorResponse := errorResponse{
			Error: "chirp is too long",
		}

		data, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		w.WriteHeader(400)
		w.Write(data)
		return
	}

	validResponse := validResponse{
		Value: true,
	}

	data, err := json.Marshal(validResponse)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	w.WriteHeader(200)
	w.Write(data)
}

type chirpRequest struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type validResponse struct {
	Value bool `json:"valid"`
}
