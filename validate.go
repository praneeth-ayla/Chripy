package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
