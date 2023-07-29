package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)


func resondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, ststusCode int, message string) {
	if ststusCode > 499 { // errors in the 400 range is the client side errors
		log.Println("Responding with 5XX error:", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	resondWithJSON(w, ststusCode, errorResponse{
		Error:	message,
	})
}
