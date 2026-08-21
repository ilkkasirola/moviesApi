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

	genreRepo := repository.NewGenreRepository(database.DB)
	genreService := service.NewGenreService(genreRepo)
	genreHandler := handlers.NewGenreHandler(genreService)

	mux := routes.SetupRoutes(movieHandler, actorHandler, genreHandler)

	// http.HandleFunc("api/movies", handler.CreateMovie)

	log.Println("server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
