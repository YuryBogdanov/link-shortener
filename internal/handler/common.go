package handler

import "net/http"

func handleError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}
