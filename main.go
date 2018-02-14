package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

// CEPSize is the size of a cep
const CEPSize = 8

// ErrorResponse is the error reponse for http requests
type ErrorResponse struct {
	Message string `json:"message"`
}

// CEP (postal address code)
type CEP struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
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

var nonDigitsRegex = regexp.MustCompile(`\D+`)

func removeNonDigits(rawCep string) string {
	return nonDigitsRegex.ReplaceAllString(rawCep, "")
}

var validPath = regexp.MustCompile("^/v1/cep/?$")

func handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Charset", "UTF-8")

	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		handleError(w, http.StatusNotFound, "Not Found")
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
		return
	}

	value, ok := query["q"]

	if !ok {
		handleError(w, http.StatusUnprocessableEntity, "CEP query missing")
		return
	}

	cep := removeNonDigits(value[0])

	if len(cep) > CEPSize {
		handleError(w, http.StatusUnprocessableEntity, "Informed CEP has more than 8 caracters")
	}

	response := CEP{
		CEP:          cep,
		Street:       "some street",
		Neighborhood: "some neighborhood",
		City:         "some city",
		State:        "some city",
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
