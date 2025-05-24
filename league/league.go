package league

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Team struct {
	Name     string
	Strength int
	Points   int
	GF       int
	GA       int
}

var Teams = []*Team{
	{"Arsenal", 80, 0, 0, 0},
	{"Chelsea", 65, 0, 0, 0},
	{"Liverpool", 90, 0, 0, 0},
	{"Manchester City", 75, 0, 0, 0},
}

var week int
var db *sql.DB

type MatchResult struct {
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeGoals int    `json:"home_goals"`
	AwayGoals int    `json:"away_goals"`
	Week      int    `json:"week"`
}

type Standing struct {
	Team   string `json:"team"`
	Points int    `json:"points"`
	GF     int    `json:"gf"`
	GA     int    `json:"ga"`
	GD     int    `json:"gd"`
}

func InitDB(filepath string) {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		week INTEGER,
		home_team TEXT,
		away_team TEXT,
		home_goals INTEGER,
		away_goals INTEGER
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("Failed to create matches table:", err)
	}
}

func SimulateWeek() []MatchResult {
	rand.Seed(time.Now().UnixNano())
	week++
	matches := [][]int{{0, 1}, {2, 3}}
	if week == 2 {
		matches = [][]int{{0, 2}, {1, 3}}
	} else if week == 3 {
		matches = [][]int{{0, 3}, {1, 2}}
	} else if week > 3 {
		return nil
	}

	results := []MatchResult{}
	for _, match := range matches {
		home := Teams[match[0]]
		away := Teams[match[1]]
		hg, ag := simulateMatch(home.Strength, away.Strength)

		home.GF += hg
		home.GA += ag
		away.GF += ag
		away.GA += hg

		if hg > ag {
			home.Points += 3
		} else if hg < ag {
			away.Points += 3
		} else {
			home.Points++
			away.Points++
		}

		// Save to DB
		_, err := db.Exec(`
			INSERT INTO matches (week, home_team, away_team, home_goals, away_goals)
			VALUES (?, ?, ?, ?, ?)`, week, home.Name, away.Name, hg, ag)
		if err != nil {
			log.Println("Error inserting match result:", err)
		}

		results = append(results, MatchResult{
			HomeTeam:  home.Name,
			AwayTeam:  away.Name,
			HomeGoals: hg,
			AwayGoals: ag,
			Week:      week,
		})
	}
	return results
}

func simulateMatch(homeStrength, awayStrength int) (int, int) {
	matchMinutes := 90
	shots := 0
	hg := 0
	ag := 0

	homePct := float64(homeStrength) / float64(homeStrength+awayStrength)

	for m := 0; m < matchMinutes; m++ {
		eventChance := rand.Intn(1000)
		if eventChance < 222 { // ~22% chance for an event (shot)
			if rand.Float64() < homePct {
				if isGoal() {
					hg++
				}
			} else {
				if isGoal() {
					ag++
				}
			}
			shots++
		}
	}
	return hg, ag
}

func isGoal() bool {
	// 1 in 5 chance that a shot results in a goal
	shotOutcome := rand.Intn(100)
	return shotOutcome < 20
}

func GetStandings() []Standing {
	standings := []Standing{}
	for _, team := range Teams {
		standings = append(standings, Standing{
			Team:   team.Name,
			Points: team.Points,
			GF:     team.GF,
			GA:     team.GA,
			GD:     team.GF - team.GA,
		})
	}
	return standings
}

func GetAllMatches() []MatchResult {
	rows, err := db.Query("SELECT week, home_team, away_team, home_goals, away_goals FROM matches ORDER BY week")
	if err != nil {
		log.Println("Error reading matches:", err)
		return nil
	}
	defer rows.Close()

	var matches []MatchResult
	for rows.Next() {
		var m MatchResult
		err := rows.Scan(&m.Week, &m.HomeTeam, &m.AwayTeam, &m.HomeGoals, &m.AwayGoals)
		if err != nil {
			log.Println("Error scanning match row:", err)
			continue
		}
		matches = append(matches, m)
	}
	return matches
}

func EditMatch(week int, homeTeam, awayTeam string, homeGoals, awayGoals int) bool {
	res, err := db.Exec(`
		UPDATE matches
		SET home_goals = ?, away_goals = ?
		WHERE week = ? AND home_team = ? AND away_team = ?
	`, homeGoals, awayGoals, week, homeTeam, awayTeam)
	if err != nil {
		log.Println("Failed to update match:", err)
		return false
	}
	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0
}

func RecalculateStandingsFromDatabase() error {
	// Reset all teams
	for _, t := range Teams {
		t.Points = 0
		t.GF = 0
		t.GA = 0
	}

	rows, err := db.Query("SELECT home_team, away_team, home_goals, away_goals FROM matches")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var homeName, awayName string
		var hg, ag int
		err := rows.Scan(&homeName, &awayName, &hg, &ag)
		if err != nil {
			return err
		}

		home := findTeamByName(homeName)
		away := findTeamByName(awayName)

		if home == nil || away == nil {
			continue
		}

		home.GF += hg
		home.GA += ag
		away.GF += ag
		away.GA += hg

		if hg > ag {
			home.Points += 3
		} else if hg < ag {
			away.Points += 3
		} else {
			home.Points++
			away.Points++
		}
	}

	return nil
}

func findTeamByName(name string) *Team {
	for _, t := range Teams {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func ResetLeague() {
	_, err := db.Exec("DELETE FROM matches")
	if err != nil {
		log.Println("Failed to reset matches table:", err)
	}
	week = 0
	for _, t := range Teams {
		t.Points = 0
		t.GF = 0
		t.GA = 0
	}
}
