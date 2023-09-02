package handlers

import "net/http"

func BadRequest(w http.ResponseWriter, r *http.Request) {
	handleBadRequest(w)
}

func handleBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}
