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
	actorRepo := repository.NewActorRepository(database.DB)
	genreRepo := repository.NewGenreRepository(database.DB)

	movieService := service.NewMovieService(movieRepo, actorRepo, genreRepo)
	actorService := service.NewActorService(actorRepo)
	genreService := service.NewGenreService(genreRepo)

	movieHandler := handlers.NewMovieHandler(movieService)
	actorHandler := handlers.NewActorHandler(actorService)
	genreHandler := handlers.NewGenreHandler(genreService)

	mux := routes.SetupRoutes(movieHandler, actorHandler, genreHandler)

	// http.HandleFunc("api/movies", handler.CreateMovie)

	log.Println("server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
