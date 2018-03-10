package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/logger"

	"fastcep/src/cache"
	"fastcep/src/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Panic(err)
	}

	cache, err := cache.NewClient(cache.GetCredentials())
	if err != nil {
		log.Panic(err)
	}
	env := &handlers.Env{DB: db, Cache: cache}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(env.SearchPostalCode)
	app := loggerMiddleware.Handler(router)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, app)
}
