package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	// connect to the SQLite database
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	err = CreateTables()
	if err != nil {
		return err
	}
	return CreateTables()
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func CreateTables() error {
	schema := `CREATE TABLE IF NOT EXISTS Genre (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS Movie (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		releaseYear INTEGER,
		duration INTEGER,
		UNIQUE(title, releaseYear, duration)
	);

	CREATE TABLE IF NOT EXISTS Actor (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		birthDate DATE,
		UNIQUE(name, birthDate)
	);

	CREATE TABLE IF NOT EXISTS MovieGenre (
	movieId INTEGER NOT NULL,
	genreId INTEGER NOT NULL,
	PRIMARY KEY (movieId, genreId),
	FOREIGN KEY (movieId) REFERENCES Movie(id) ON DELETE CASCADE,
	FOREIGN KEY (genreId) REFERENCES Genre(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS MovieActor (
	movieId INTEGER NOT NULL,
	actorId INTEGER NOT NULL,
	PRIMARY KEY (movieId, actorId),
	FOREIGN KEY (movieId) REFERENCES Movie(id) ON DELETE CASCADE,
	FOREIGN KEY (actorId) REFERENCES Actor(id) ON DELETE CASCADE
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Println("Failed to create table:", err)
		return err
	}
	log.Println("All tables created/verified successfully")
	return err
}
