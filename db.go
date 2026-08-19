package main

import (
	"database/sql"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
)

type Actor struct {
	id        int
	name      string
	birthDate int
}

func CreateTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Actor (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		birthDate TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println("ERROR")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Genre (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
		)`)
	if err != nil {
		fmt.Println("error")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Movie (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			releaseYear INTEGER NOT NULL,
			duration INTEGER NOT NULL
		)`)
	if err != nil {
		fmt.Println("error")
	}
}
