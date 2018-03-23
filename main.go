package main

import (
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/logger"

	"fastcep/src/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	env := &handlers.Env{DB: db}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(env.SearchPostalCode)
	app := loggerMiddleware.Handler(router)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, app)
}
