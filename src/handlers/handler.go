package handlers

import (
	"database/sql"
	"encoding/json"
	"fastcep/src/cep"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/go-redis/redis"
)

// Env holds the environment connections
type Env struct {
	DB    *sql.DB
	Cache *redis.Client
}

var validPath = regexp.MustCompile("^/v1/cep/?$")

// SearchPostalCode on the database registered data
func (env *Env) SearchPostalCode(w http.ResponseWriter, r *http.Request) {
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

	cepValue := cep.RemoveNonDigits(value[0])

	if len(cepValue) > cep.CEPSize {
		handleError(w, http.StatusUnprocessableEntity, "Informed CEP has more than 8 caracters")
		return
	}

	cepValue = cep.LeftPadZero(cepValue, cep.CEPSize)

	val, err := env.Cache.Get("cep:" + cepValue).Result()
	if err == redis.Nil {
		var response cep.CEP
		row := env.DB.QueryRow("SELECT p.cep, p.street, p.neighborhood, s.name AS state, c.name AS city FROM postal_codes AS p INNER JOIN states AS s ON s.id=p.state_id INNER JOIN cities AS c ON c.id=p.city_id WHERE p.cep=$1", cepValue)

		err = row.Scan(&response.CEP, &response.Street, &response.Neighborhood, &response.State, &response.City)

		switch {
		case err == sql.ErrNoRows:
			message := fmt.Sprintf("CEP número %s não foi encontrado", cepValue)
			handleError(w, http.StatusNotFound, message)
		case err != nil:
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
		default:
			content, err := json.Marshal(response)

			if err != nil {
				handleError(w, http.StatusInternalServerError, "Internal  Server Error")
				return
			}

			err = env.Cache.Set("cep:"+cepValue, content, 6*time.Hour).Err()
			if err != nil {
				handleError(w, http.StatusInternalServerError, "Internal  Server Error")
				return
			}

			var cepResponse cep.CEP
			err = json.Unmarshal(content, &cepResponse)

			if err != nil {
				handleError(w, http.StatusInternalServerError, "Internal  Server Error")
				return
			}

			err = json.NewEncoder(w).Encode(cepResponse)

			if err != nil {
				handleError(w, http.StatusInternalServerError, "Internal  Server Error")
				return
			}
		}
	} else if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
		return
	} else {
		var cepResponse cep.CEP
		err = json.Unmarshal([]byte(val), &cepResponse)

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
			return
		}

		err = json.NewEncoder(w).Encode(cepResponse)

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
			return
		}
	}
}
