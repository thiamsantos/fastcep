package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/logger"

	"fastcep/src/handlers"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Panic(err)
	}

	env := &handlers.Env{DB: db}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(env.SearchPostalCode)
	app := loggerMiddleware.Handler(router)
	http.ListenAndServe(":8080", app)
}
