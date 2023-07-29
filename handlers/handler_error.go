package handlers

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
