package handlers

import (
	"database/sql"
	"encoding/json"
	"fastcep/src/address"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	cache "github.com/patrickmn/go-cache"
)

// Env holds the environment connections
type Env struct {
	DB    *sql.DB
	Cache *cache.Cache
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

	cepValue := address.RemoveNonDigits(value[0])

	if len(cepValue) > address.CEPSize {
		handleError(w, http.StatusUnprocessableEntity, "Informed CEP has more than 8 caracters")
		return
	}

	cepValue = address.LeftPadZero(cepValue, address.CEPSize)

	val, found := env.Cache.Get("cep:" + cepValue)

	if found {
		err = json.NewEncoder(w).Encode(val.(address.Address))

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
			return
		}
		return
	}

	var response address.Address
	row := env.DB.QueryRow("SELECT p.cep, p.street, p.neighborhood, p.state, p.city, p.uf FROM postal_codes AS p WHERE p.cep=$1", cepValue)

	err = row.Scan(&response.CEP, &response.Street, &response.Neighborhood, &response.State, &response.City, &response.Uf)

	switch {
	case err == sql.ErrNoRows:
		message := fmt.Sprintf("CEP número %s não foi encontrado", cepValue)
		handleError(w, http.StatusNotFound, message)
	case err != nil:
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
	default:
		env.Cache.Set("cep:"+cepValue, response, cache.DefaultExpiration)

		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			handleError(w, http.StatusInternalServerError, "Internal  Server Error")
			return
		}
	}

}
