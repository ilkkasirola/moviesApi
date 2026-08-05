package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./movie.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to the SQLite database.")

	CreateTable(db)

	fmt.Println("Tables created.")
}
