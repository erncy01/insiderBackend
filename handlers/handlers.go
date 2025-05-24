package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"insiderBackend/league"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "🏟️ Insider League Simulator is running!")
}

func SimulateWeekHandler(w http.ResponseWriter, r *http.Request) {
	results := league.SimulateWeek()
	json.NewEncoder(w).Encode(results)
}

func SimulateAllHandler(w http.ResponseWriter, r *http.Request) {
	var allResults [][]league.MatchResult
	for i := 0; i < 3; i++ {
		results := league.SimulateWeek()
		if results == nil {
			break
		}
		allResults = append(allResults, results)
	}
	json.NewEncoder(w).Encode(allResults)
}

func StandingsHandler(w http.ResponseWriter, r *http.Request) {
	standings := league.GetStandings()
	json.NewEncoder(w).Encode(standings)
}

func MatchesHandler(w http.ResponseWriter, r *http.Request) {
	matches := league.GetAllMatches()
	json.NewEncoder(w).Encode(matches)
}

func EditMatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type EditInput struct {
		Week      int    `json:"week"`
		HomeTeam  string `json:"home_team"`
		AwayTeam  string `json:"away_team"`
		HomeGoals int    `json:"home_goals"`
		AwayGoals int    `json:"away_goals"`
	}
	var input EditInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	success := league.EditMatch(input.Week, input.HomeTeam, input.AwayTeam, input.HomeGoals, input.AwayGoals)
	if success {
		err := league.RecalculateStandingsFromDatabase()
		if err != nil {
			http.Error(w, "Match updated but failed to recalculate standings", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "Match updated and standings recalculated successfully")
	} else {
		http.Error(w, "Match not found", http.StatusNotFound)
	}
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	league.ResetLeague()
	fmt.Fprint(w, "League has been reset.")
}
