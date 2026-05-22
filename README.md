# Go Football League

A football league simulation REST API built with **Go**, **Gin**, and **PostgreSQL**.

Simulates a full 20-team Premier League season (38 weeks, 380 fixtures) using real EA FC 26 team ratings (Attack, Defense, Midfield). Match scores are generated based on a  Poisson model to adapt real Premier League goal averages. Championship probabilities are calculated from week 4.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Framework | Gin |
| Database | PostgreSQL |
| DB Driver | sqlx + lib/pq |
| Config | godotenv |

---

## Architecture

Strict layered architecture — each layer depends only on the interface of the layer below it (dependency injection throughout).

```
cmd/main.go             → wires all layers, seeds fixtures, starts server
internal/
  handler/              → HTTP layer (Gin context, request binding, JSON responses)
  service/              → business logic (simulation, standings, predictions)
  repository/           → database queries (sqlx + PostgreSQL)
  models/               → domain structs + interfaces shared across layers
pkg/
  database/             → DB connection (godotenv) + fixture seeder
migrations/
  001_init.sql          → schema (teams, matches) + 20 Premier League teams
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL running locally

### 1. Clone the repo

```bash
git clone https://github.com/dogancanbilgili/Go_Football_League.git
cd Go_Football_League
```

### 2. Set up the database

Create a PostgreSQL database:

```sql
CREATE DATABASE go_football_league;
```

Run the migration to create tables and insert the 20 teams:

```bash
psql -U postgres -d go_football_league -f migrations/001_init.sql
```

### 3. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=go_football_league
```

### 4. Install dependencies and run

```bash
go mod download
go run ./cmd
```

The server starts on `http://localhost:8080`. On first run, the seeder automatically generates all 380 fixtures using the circle-method round-robin algorithm and inserts them in a single transaction.

### Live deployment

The API is deployed and publicly accessible at:

```
https://gofootballleague-production.up.railway.app
```

The frontend is live at:

```
https://dogancanbilgili.github.io/Go_Football_League/
```

---

## API Reference

All endpoints are under the `/league` prefix. Test with Postman or curl.

### GET `/league/table`

Returns the current league standings sorted by points → goal difference → goals for → alphabetical.

```json
[
  {
    "id": 7, "name": "Manchester City",
    "attack": 88, "midfield": 86, "defence": 84,
    "played": 5, "won": 4, "drawn": 0, "lost": 1,
    "goals_for": 12, "goals_against": 4, "goal_difference": 8, "points": 12
  }
]
```

---

### GET `/league/fixtures`

Returns all 380 fixtures. Unplayed matches have `null` for goals and `"played": false`.

---

### GET `/league/week/:week`

Returns the 10 fixtures (and results if played) for a specific week.

```
GET /league/week/3
```

---

### POST `/league/next-week`

Simulates the next unplayed week. Returns the 10 match results.

Returns `400` if all 38 weeks have already been played.

---

### POST `/league/play-all`

Simulates all remaining weeks at once. Returns results grouped by week number.

```json
{
  "6":  [ { "week": 6, "home_team": { "name": "Arsenal" }, "away_team": { "name": "Chelsea" }, "home_goals": 2, "away_goals": 1 } ],
  "38": [ ]
}
```

---

### PUT `/league/match/:id`

Manually override a match result. The league table recalculates automatically on the next `GET /league/table` call — standings are computed from match data, not stored separately.

```json
{
  "home_goals": 3,
  "away_goals": 1
}
```

Returns `400` if the match ID is invalid or goals are negative.

---

### GET `/league/predictions`

Returns championship probability for each team based on current points plus expected points from remaining fixtures.

Available from **week 4 onward**. Returns `400` before that.

```json
[
  { "name": "Manchester City", "percentage": 34.21 },
  { "name": "Arsenal",         "percentage": 28.57 }
]
```

---

### POST `/league/reset`

Resets all match results to `null` / `played = false`. Teams and fixture schedule are preserved.

```json
{ "message": "league reset to initial state" }
```

---

## Simulation Model

### Match simulation — Poisson distribution

Each team's expected goal count (lambda) is derived from the attack vs. defence matchup:

```
lambdaHome = (homeAttack / (homeAttack + awayDefence)) × 2.6
lambdaAway = (awayAttack / (awayAttack + homeDefence)) × 2.2
```

The multipliers (2.6 / 2.2) are calibrated to real Premier League averages (~1.5 home goals, ~1.2 away goals per game). The gap between them models home advantage.

Goals are drawn using **Knuth's algorithm**: multiply uniform random numbers until their product falls below `e^(-lambda)`. The number of draws required follows a Poisson distribution with mean = lambda.

### Championship prediction

From week 4 onward, each team's championship score is:

```
score = currentPoints + Σ (expected points from each remaining fixture)
```

Expected points per fixture uses the team's attack + midfield rating relative to the combined contest, weighted by home/away probability (55% / 45%). Each team's percentage is `score / totalScore × 100`.

### Team ratings

20 Premier League teams with real **EA FC 26** ratings sourced from [sofifa.com](https://sofifa.com). Three attributes per team: Attack, Midfield, Defence. Attack and Defence drive match simulation; Midfield drives championship predictions.

---

## Database Schema

```sql
CREATE TABLE teams (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(100) NOT NULL UNIQUE,
    attack   INT NOT NULL,
    midfield INT NOT NULL,
    defence  INT NOT NULL
);

CREATE TABLE matches (
    id           SERIAL PRIMARY KEY,
    week         INT NOT NULL,
    home_team_id INT NOT NULL REFERENCES teams(id),
    away_team_id INT NOT NULL REFERENCES teams(id),
    home_goals   INT DEFAULT NULL,
    away_goals   INT DEFAULT NULL,
    played       BOOLEAN NOT NULL DEFAULT false
);
```

Standings are **computed at query time** from match results — not stored. This means editing a match result is immediately reflected in the table with no additional writes.

---

## Project Structure

```
.
├── cmd/
│   └── main.go                    # entry point — wires layers, starts server
├── internal/
│   ├── handler/
│   │   └── league_handler.go      # all 8 HTTP endpoints
│   ├── service/
│   │   └── league_service.go      # Poisson simulation, table, predictions
│   ├── repository/
│   │   ├── team_repository.go     # team DB queries
│   │   └── match_repository.go    # match DB queries
│   └── models/
│       ├── team.go                # Team, Standing, Prediction structs
│       ├── match.go               # Match struct
│       ├── repository.go          # TeamRepository, MatchRepository interfaces
│       └── service.go             # LeagueService interface
├── pkg/
│   └── database/
│       ├── postgres.go            # DB connection via godotenv
│       └── seed.go                # circle-method fixture generator
├── migrations/
│   └── 001_init.sql               # schema + 20 team inserts
├── .env.example
├── go.mod
└── README.md
```
