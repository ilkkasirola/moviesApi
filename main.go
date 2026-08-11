package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"log"
	"movies-api/database"
	"movies-api/routes"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	if err := database.InitDB("./movies.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	mux := routes.SetupRoutes()

	log.Println("serer on http://localhost:8080")
	insertTestData()
	log.Fatal(http.ListenAndServe(":8080", mux))

	var sqliteVersion string
	err := database.DB.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println("Error getting SQLite version:", err)
		return
	}
	fmt.Println("SQLite Version:", sqliteVersion)

	listTables(database.DB)

}

func insertTestData() {
	// Insert genres
	database.DB.Exec("INSERT OR IGNORE INTO Genre (name) VALUES ('Action')")
	database.DB.Exec("INSERT OR IGNORE INTO Genre (name) VALUES ('Thriller')")

	// Insert actors
	database.DB.Exec("INSERT OR IGNORE INTO Actor (name, birthDate) VALUES ('Bruce Willis', '1955-03-19')")
	database.DB.Exec("INSERT OR IGNORE INTO Actor (name, birthDate) VALUES ('Alan Rickman', '1946-02-21')")

	// Insert movie
	database.DB.Exec("INSERT OR IGNORE INTO Movie (title, releaseYear, duration) VALUES ('Die Hard', 1988, 131)")

	// Link genres (assuming ID 1 = Die Hard, 1 = Action, 2 = Thriller)
	database.DB.Exec("INSERT OR IGNORE INTO MovieGenre (movieId, genreId) VALUES (1, 1)")
	database.DB.Exec("INSERT OR IGNORE INTO MovieGenre (movieId, genreId) VALUES (1, 2)")

	// Link actors (assuming ID 1 = Bruce Willis, 2 = Alan Rickman)
	database.DB.Exec("INSERT OR IGNORE INTO MovieActor (movieId, actorId) VALUES (1, 1)")
	database.DB.Exec("INSERT OR IGNORE INTO MovieActor (movieId, actorId) VALUES (1, 2)")
}

func listTables(db *sql.DB) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;")
	if err != nil {
		log.Fatal("Error querying tables:", err)
	}
	defer rows.Close()

	fmt.Println("\nTables in database:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal("Error scanning table name:", err)
		}
		fmt.Println(" -", tableName)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error during rows iteration", err)
	}
}
