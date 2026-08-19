package main

import (
	"net/http"

	"log"
	"movies-api/database"
	"movies-api/handlers"
	"movies-api/repository"
	"movies-api/routes"
	"movies-api/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := database.InitDB("./movies.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	repo := repository.NewMovieRepository(database.DB)
	service := service.NewMovieService(repo)
	handler := handlers.NewMovieHandler(service)

	mux := routes.SetupRoutes(handler)

	http.HandleFunc("api/movies", handler.CreateMovie)

	log.Println("serer on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
