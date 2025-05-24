// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"insider/handlers"
	"insider/league"
)

func main() {
	league.InitDB("league.db")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ping", handlers.PingHandler)
	http.HandleFunc("/simulate-week", handlers.SimulateWeekHandler)
	http.HandleFunc("/simulate-all", handlers.SimulateAllHandler)
	http.HandleFunc("/standings", handlers.StandingsHandler)
	http.HandleFunc("/matches", handlers.MatchesHandler)
	http.HandleFunc("/edit-match", handlers.EditMatchHandler)
	http.HandleFunc("/reset", handlers.ResetHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
