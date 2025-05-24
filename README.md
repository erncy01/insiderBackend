# Insider League Simulator (GoLang)

This project simulates a mini football league with 4 teams using GoLang for backend and SQLite for persistence. It includes a lightweight HTML/JS frontend served by Go.

## Features

- Simulate one week or an entire league season
- View real-time league standings and match results
- Edit past match results and automatically recalculate standings
- All data stored in a SQLite database
- Clean, single-page frontend with no external libraries or React

## Setup Instructions (Local)

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- [GCC installed and in your](https://sourceforge.net/projects/tdm-gcc/)
- SQLite C compiler support (CGO-enabled)

### Step-by-Step

1- Clone the repository

`git clone https://github.com/erncy01/insiderBackend.git`

2- Navigate into the folder.

`cd insiderBackend`

3- Enable CGO to use SQLite

`set CGO_ENABLED=1`

4- Run the program

`go run main.go`

5- Visit the local port
`(http://localhost:8080)`


---

## 🧠 How Simulation Works

- [Based on this the Match Simulator algorithm.](https://matchsimulator.com/about)
- Teams have a base **strength** (e.g., 80 vs 75).
- Match is simulated **minute-by-minute (90 mins)**.
- Each minute has a **22% chance for a shot**.
- Stronger team more likely to take the shot (based on % share).
- Each shot has a **20% chance of being a goal**.
- This mimics Poisson-style football simulation.

---

## 📄 File-by-File Explanation

### `main.go`
- Entry point.
- Initializes DB.
- Serves static files from `/static/`.
- Routes API endpoints (`/simulate-week`, `/standings`, etc.).

### `handlers/handlers.go`
- HTTP handler functions.
- Bridges Go logic (`league/`) with HTTP API.
- Serializes JSON responses for frontend.

### `league/league.go`
- Core league mechanics and DB access.
- Functions:
  - `SimulateWeek()`: Realistic strength-based match simulation.
  - `RecalculateStandingsFromDatabase()`: Rebuilds team stats.
  - `EditMatch()`: Updates DB and recalculates.
  - `ResetLeague()`: Clears all state.
  - `GetAllMatches()` & `GetStandings()` for frontend.

### `static/index.html`
- One-page frontend UI.
- Buttons for simulation/reset.
- Match editor form.
- Tables for standings and history.

### `static/script.js`
- Handles all button actions via `fetch`.
- Sends and receives JSON.
- Renders standings/matches in HTML tables.
### HTML and JavaScript only used to make the frontend look better.
---

## API Endpoints

| Endpoint          | Method | Description                           |
|------------------|--------|---------------------------------------|
| `/ping`          | GET    | Server test                           |
| `/simulate-week` | GET    | Simulates 1 week                      |
| `/simulate-all`  | GET    | Simulates 3 weeks                     |
| `/standings`     | GET    | JSON of team standings                |
| `/matches`       | GET    | JSON of all match results             |
| `/edit-match`    | PUT    | Update a match result in DB           |
| `/reset`         | GET    | Clears DB and resets teams            |

## SQL Scheme and Queries
SQL database is created in Go code. Scheme and queries used in the code with their goal is given below.

### Scheme
| Column      | Data Type | Description                         | Notes            |
| ----------- | --------- | ----------------------------------- | ---------------- |
| id          | INTEGER   | Primary key                         | Auto-incremented |
| week        | INTEGER   | Match week number                   |                  |
| home\_team  | TEXT      | Name of the home team               |                  |
| away\_team  | TEXT      | Name of the away team               |                  |
| home\_goals | INTEGER   | Number of goals scored by home team |                  |
| away\_goals | INTEGER   | Number of goals scored by away team |                  |

### Queries

### 1. Creating the matches table
```sql
CREATE TABLE IF NOT EXISTS matches (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  week INTEGER,
  home_team TEXT,
  away_team TEXT,
  home_goals INTEGER,
  away_goals INTEGER
);
```
- Creates the matches table with columns for match details.
- Ensures the table is created only if it does not exist.

### 2. Inserting a new match result
```sql
INSERT INTO matches (week, home_team, away_team, home_goals, away_goals)
VALUES (?, ?, ?, ?, ?);
```
- Inserts a new match result into the database.
- Uses placeholders (?) for dynamic values passed from the Go code.

### 3. Selecting all matches ordered by week
```sql
SELECT week, home_team, away_team, home_goals, away_goals FROM matches ORDER BY week;
```
- Retrieves all matches, ordered by the week number for sequential display or processing.

### 4. Updating an existing match result
```sql
UPDATE matches
SET home_goals = ?, away_goals = ?
WHERE week = ? AND home_team = ? AND away_team = ?;
```
- Updates the goals of a specific match identified by week and teams.
- Allows correcting or editing match results after they are recorded.

### 5. Selecting matches for recalculating standings
```sql
SELECT home_team, away_team, home_goals, away_goals FROM matches;
```
- Fetches all matches to recalculate points, goals for/against, and standings for all teams.

### 6. Deleting all matches to reset the league
```sql
DELETE FROM matches;
```
- Deletes all match records to reset the league data.
- Used when restarting the simulation or clearing results.


