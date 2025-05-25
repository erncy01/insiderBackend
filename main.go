// Package main initializes and starts the server.
package main

import (
	"fmt"
	"log"
	"net/http"

	"insiderBackend/handlers"
	"insiderBackend/league"
)

func main() {
	// Initialize the repository for match data storage.
	league.InitDB("league.db")

	// Serve static files from the "static" directory.
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Define HTTP routes and associate them with handler functions.
	http.HandleFunc("/ping", handlers.PingHandler)
	http.HandleFunc("/simulate-week", handlers.SimulateWeekHandler)
	http.HandleFunc("/simulate-all", handlers.SimulateAllHandler)
	http.HandleFunc("/standings", handlers.StandingsHandler)
	http.HandleFunc("/matches", handlers.MatchesHandler)
	http.HandleFunc("/edit-match", handlers.EditMatchHandler)
	http.HandleFunc("/reset", handlers.ResetHandler)

	// Start the HTTP server on port 8080.
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
