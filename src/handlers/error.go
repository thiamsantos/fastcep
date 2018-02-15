package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the error reponse for http requests
type ErrorResponse struct {
	Message string `json:"message"`
}

func handleError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	response := ErrorResponse{message}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
