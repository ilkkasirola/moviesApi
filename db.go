package main

import (
	"database/sql"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
)

type Actor struct {
	id        int
	name      string
	birthDate string
}

type Genre struct {
	id   int
	name string
}

type Movie struct {
	id          int
	title       string
	releaseYear int
	duration    int
	Actor       Actor
	Genre       Genre
}

func CreateTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Actors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		birthDate TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS Genres (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
		);

	CREATE TABLE IF NOT EXISTS Movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			releaseYear INTEGER NOT NULL,
			duration INTEGER NOT NULL
			FOREIGN KEY(Genres_id) REFERENCES Genres(id)
			FOREIGN KEY(Actors_id) REFERENCES Actors(id)
		)`)
	if err != nil {
		fmt.Println("Error creating tables.")
	}
}
