package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/unrolled/logger"

	"fastcep/src/database"
	"fastcep/src/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewPool(database.GetCredentials())
	if err != nil {
		log.Panic(err)
	}
	env := &handlers.Env{DB: db}

	loggerMiddleware := logger.New()
	router := http.HandlerFunc(env.QueryPostalCode)
	app := loggerMiddleware.Handler(router)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, app)
}
