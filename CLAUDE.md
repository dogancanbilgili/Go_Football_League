# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

A football league simulation REST API (no frontend). Four teams play a full round-robin season (each pair plays home and away = 6 weeks, 2 matches per week). The API exposes endpoints consumable via Postman. The task is from an Insider internship hiring challenge.

## Commands

```bash
go mod download       # install dependencies
go run ./cmd          # start the server
go build ./cmd        # build binary
go test ./...         # run all tests
go vet ./...          # static analysis
```

Environment: copy `.env.example` to `.env` and fill in database credentials before running.

## Architecture

Strict layered architecture under `internal/`; nothing in `internal/` is importable from outside the module.

```
cmd/main.go           → wires dependencies, starts Gin server
internal/handler/     → HTTP handlers (Gin context, request binding, response)
internal/service/     → business logic (league simulation, championship prediction)
internal/repository/  → database queries via sqlx + lib/pq (PostgreSQL)
internal/models/      → domain structs and interfaces
pkg/database/         → DB connection setup (godotenv + lib/pq)
migrations/           → SQL schema and seed files
```

**Required design constraint:** use Go interfaces and struct composition throughout. Each layer depends on an interface from the layer below, not a concrete type. Example pattern:

```go
// models defines the contract
type LeagueRepository interface { ... }

// repository implements it
type leagueRepo struct { db *sqlx.DB }

// service depends on the interface, not the concrete type
type leagueService struct { repo models.LeagueRepository }
```

## Domain Rules

**League:** 4 teams (e.g. Chelsea, Arsenal, Manchester City, Liverpool), each with a `strength` attribute (1–100) that biases match simulation. Teams play each other home and away → 6 weeks total, 2 fixtures per week.

**Scoring (Premier League rules):** Win = 3 pts, Draw = 1 pt, Loss = 0 pts. Table sorted by: points → goal difference → goals for → alphabetical.

**Championship prediction:** shown starting from week 4 onward. Each team's probability is proportional to its current points + remaining strength-weighted expected points.

**Match simulation:** goals scored per team modelled from team strength (e.g. Poisson distribution with λ derived from attacker strength vs. opponent defender strength).

## Key API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/league/table` | Current league standings |
| GET | `/league/fixtures` | All fixtures with results |
| POST | `/league/next-week` | Simulate the next week's matches |
| POST | `/league/play-all` | Simulate all remaining weeks at once |
| GET | `/league/week/:week` | Match results for a specific week |
| PUT | `/league/match/:id` | Edit a match result and recalculate table |
| GET | `/league/predictions` | Championship probability percentages (week ≥ 4) |
| POST | `/league/reset` | Reset the league to initial state |

## Database

Primary: **PostgreSQL** via `lib/pq` + `sqlx`. SQL schema lives in `migrations/`. Core tables: `teams`, `matches`. The `matches` table stores home/away team IDs, scores, week number, and a `played` boolean.

MongoDB driver is in go.mod but PostgreSQL is the primary store; do not introduce dual-write complexity unless explicitly required.

## Module

`github.com/dogancanbilgili/Go_Football_League` — Go 1.26.3
