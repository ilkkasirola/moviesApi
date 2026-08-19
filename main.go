package main

import (
	"net/http"

	"log"
	"movies-api/database"
	"movies-api/routes"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
<<<<<<< HEAD

	db, err := sql.Open("sqlite3", "./movie.db?_foreign_keys=ON")
	if err != nil {
		fmt.Println(err)
		return
=======
	if err := database.InitDB("./movies.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
>>>>>>> cde7176c5d58bca99acf59bc63c4fdceddef2870
	}
	defer database.CloseDB()

	mux := routes.SetupRoutes()

	log.Println("serer on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

<<<<<<< HEAD
	fmt.Println("Tables created.")
	RunServer(db)
=======
>>>>>>> cde7176c5d58bca99acf59bc63c4fdceddef2870
}
