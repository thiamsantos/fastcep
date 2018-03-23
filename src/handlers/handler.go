package handlers

import (
	"fastcep/src/address"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/boltdb/bolt"
)

// Env holds the environment connections
type Env struct {
	DB *bolt.DB
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
	err = env.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("postal_codes"))
		v := b.Get([]byte(cepValue))

		if v == nil {
			message := fmt.Sprintf("CEP número %s não foi encontrado", cepValue)
			handleError(w, http.StatusNotFound, message)
			return nil
		}

		w.Write(v)

		return nil
	})

	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal  Server Error")
		return
	}
}
