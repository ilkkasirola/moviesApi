package main

import (
	"net/http"

	"log"
	"movies-api/database"
	"movies-api/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := database.InitDB("./movies.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	mux := routes.SetupRoutes()

	log.Println("server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
