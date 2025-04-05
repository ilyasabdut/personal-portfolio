package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./sqlite/data.db")
	if err != nil {
		log.Fatalf("failed to connect to sqlite: %v", err)
	}

	log.Println("Connected to SQLite database")
	
	createTable := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		tech_stack TEXT NOT NULL,
		deployed TEXT NOT NULL,
		description TEXT
	);
	`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("failed to create projects table: %v", err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
