package main

import (
	"log"
	"net/http"

	"portfolio/database"
	"portfolio/handlers"
)

func main() {
	database.InitDB() 
	defer database.CloseDB()

	// Static files (optional: if you add CSS, JS, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Routes
	http.HandleFunc("/", handlers.HomePageHandler)
	http.HandleFunc("/search", handlers.SearchProjectsHandler)
	http.HandleFunc("/chat", handlers.ChatHandler)



	log.Println("Running server at http://0.0.0.0:8081")
	http.ListenAndServe("0.0.0.0:8081", nil)
}
