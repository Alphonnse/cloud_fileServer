package handlers

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, 400, "Something went wrong")
}
