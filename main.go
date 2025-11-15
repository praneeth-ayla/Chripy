package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var cfg apiConfig

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Body string `json:"body"`
	}

	type ResBody struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	reqBody := ReqBody{}
	resBody := ResBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding body: %s", err)
		w.WriteHeader(500)
		resBody.Error = "Something went wrong"
		return
	}

	if len(reqBody.Body) > 140 {
		w.WriteHeader(400)
		resBody.Error = "Chirp is too long"
		resBody.Valid = false
	} else {
		w.WriteHeader(200)
		resBody.Valid = true
	}

	dat, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		resBody.Error = "Something went wrong"
		return
	}

	w.Write(dat)
}
