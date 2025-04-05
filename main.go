package main

import (
	"fmt"
	"log"
	"net/http"

	"portfolio/database"
	"portfolio/handlers"
)

func main() {
	database.InitDB() 
	defer database.CloseDB()

	// Routes
	http.HandleFunc("/", handlers.HomePageHandler)
	http.HandleFunc("/search", handlers.SearchProjectsHandler)

	// Static files (optional: if you add CSS, JS, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


	fmt.Println("Running server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
