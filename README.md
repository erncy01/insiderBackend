# Insider League Simulator (GoLang + HTML)

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

