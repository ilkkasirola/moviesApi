package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	// connect to the SQLite database
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return err
	}
	_, err = DB.Exec("PRAGMA foreign_keys = on")

	// defer DB.Close()
	// log.Println("Connected to the SQLite database successfully.")

	// get the version of SQLite
	err = CreateTables()
	if err != nil {

	}
	return nil
}

func CreateTables() error {
	sqlGenreTable := `CREATE TABLE IF NOT EXISTS Genre (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	`

	sqlMovieTable := `
	CREATE TABLE IF NOT EXISTS Movie (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		releaseYear INTEGER,
		duration INTEGER
	);
	`

	sqlActorTable := `
	CREATE TABLE IF NOT EXISTS Actor (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		birthDate DATE
	);
	`

	sqlMovieGenreJunction := `
	CREATE TABLE IF NOT EXISTS MovieGenre (
	movieId INTEGER NOT NULL,
	genreId INTEGER NOT NULL,
	PRIMARY KEY (movieId, GenreId),
	FOREIGN KEY (movieId) REFERENCES Movie(id) ON DELETE CASCADE,
	FOREIGN KEY (genreId) REFERENCES Genre(id) ON DELETE CASCADE
	);
	`
	sqlMovieActorJunction := `
	CREATE TABLE IF NOT EXISTS MovieActor (
	movieId INTEGER NOT NULL,
	actorId INTEGER NOT NULL,
	PRIMARY KEY (movieId, actorId),
	FOREIGN KEY (movieId) REFERENCES Movie(id) ON DELETE CASCADE,
	FOREIGN KEY (actorId) REFERENCES Actor(id) ON DELETE CASCADE
	);
	`

	tables := []string{sqlGenreTable, sqlMovieTable, sqlActorTable, sqlMovieGenreJunction, sqlMovieActorJunction}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatal("Failed to create table:", err)
			return err
		}
	}
	log.Println("All tables created/verified successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
