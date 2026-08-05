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
		log.Fatal("Failed to open database:", err)
		return err
	}

	// defer DB.Close()
	// log.Println("Connected to the SQLite database successfully.")

	// get the version of SQLite
	err = CreateTables()
	return nil
}

func CreateTables() error {
	sqlGenreTable := `CREATE TABLE IF NOT EXISTS Genre (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
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
		name TEXT NOT NULL,
		birthDate DATE
	);
	`
	tables := []string{sqlGenreTable, sqlMovieTable, sqlActorTable}

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
