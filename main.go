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
		log.Fatal(err)
	}
	defer database.CloseDB()

	movieRepo := repository.NewMovieRepository(database.DB)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	actorRepo := repository.NewActorRepository(database.DB)
	actorService := service.NewActorService(actorRepo)
	actorHandler := handlers.NewActorHandler(actorService)

	mux := routes.SetupRoutes(movieHandler, actorHandler)

	// http.HandleFunc("api/movies", handler.CreateMovie)

	log.Println("serer on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
