package main

import (
	"database/sql"
	"fmt"
	"log"

	"movies-api/database"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	if err := database.InitDB("./movies.db"); err != nil {
		log.Fatal("Failed to  initialize database:", err)
	}
	defer database.CloseDB()

	var sqliteVersion string
	err := database.DB.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println("Error getting SQLite version:", err)
		return
	}
	fmt.Println("SQLite Version:", sqliteVersion)

	listTables(database.DB)
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
