package repository

import (
	"github.com/dogancanbilgili/Go_Football_League/internal/models"
	"github.com/jmoiron/sqlx"
)

type matchRepo struct {
	db *sqlx.DB
}

// NewMatchRepository returns a MatchRepository backed by PostgreSQL.
func NewMatchRepository(db *sqlx.DB) models.MatchRepository {
	return &matchRepo{db: db}
}

func (r *matchRepo) GetAll() ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Select(&matches, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, played
		FROM matches
		ORDER BY week, id
	`)
	return matches, err
}

func (r *matchRepo) GetByWeek(week int) ([]models.Match, error) {
	var matches []models.Match
	err := r.db.Select(&matches, `
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals, played
		FROM matches
		WHERE week = $1
		ORDER BY id
	`, week)
	return matches, err
}

// GetCurrentWeek returns the next unplayed week number.
// If all matches are played it returns 7 (beyond the last week).
func (r *matchRepo) GetCurrentWeek() (int, error) {
	var week int
	err := r.db.Get(&week, `
		SELECT COALESCE(MIN(week), 39)
		FROM matches
		WHERE played = false
	`)
	return week, err
}

func (r *matchRepo) UpdateMatch(id, homeGoals, awayGoals int) error {
	_, err := r.db.Exec(`
		UPDATE matches
		SET home_goals = $1, away_goals = $2, played = true
		WHERE id = $3
	`, homeGoals, awayGoals, id)
	return err
}

func (r *matchRepo) SaveSimulatedMatches(matches []models.Match) error {
	for _, m := range matches {
		if err := r.UpdateMatch(m.ID, *m.HomeGoals, *m.AwayGoals); err != nil {
			return err
		}
	}
	return nil
}

func (r *matchRepo) ResetMatches() error {
	_, err := r.db.Exec(`
		UPDATE matches
		SET home_goals = NULL, away_goals = NULL, played = false
	`)
	return err
}
