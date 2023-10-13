package handlers

import "net/http"

type response struct {
	State string `json:"state"`
}

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	ResondWithJSON(w, 200, response{
		"Server works",
	})
}
