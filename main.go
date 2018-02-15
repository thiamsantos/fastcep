package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
	"github.com/unrolled/logger"

	"github.com/joho/godotenv"
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

func leftPadZero(rawCep string, length int) string {
	return strings.Repeat("0", length-len(rawCep)) + rawCep
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

	cep = leftPadZero(cep, CEPSize)

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/fastcep?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	var response CEP
	row := db.QueryRow("SELECT p.cep, p.street, p.neighborhood, s.name AS state, c.name AS city FROM postal_codes AS p INNER JOIN states AS s ON s.id=p.state_id INNER JOIN cities AS c ON c.id=p.city_id WHERE p.cep=$1", cep)

	err = row.Scan(&response.CEP, &response.Street, &response.Neighborhood, &response.State, &response.City)

	switch {
	case err == sql.ErrNoRows:
		message := fmt.Sprintf("CEP número %s não foi encontrado", cep)
		handleError(w, http.StatusNotFound, message)
	case err != nil:
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
	default:
		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
			return
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(handler)
	app := loggerMiddleware.Handler(router)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, app)
}
