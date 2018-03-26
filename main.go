package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
	"github.com/unrolled/logger"

	"fastcep/src/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Panic(err)
	}

	c := cache.New(5*time.Minute, 10*time.Minute)

	env := &handlers.Env{DB: db, Cache: c}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(env.SearchPostalCode)
	app := loggerMiddleware.Handler(router)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, app)
}
