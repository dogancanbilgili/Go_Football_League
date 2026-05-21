package repository

import (
	"github.com/dogancanbilgili/Go_Football_League/internal/models"
	"github.com/jmoiron/sqlx"
)

type teamRepo struct {
	db *sqlx.DB
}

// NewTeamRepository returns a TeamRepository backed by PostgreSQL.
// It returns the interface type, not the concrete struct — callers never see teamRepo directly.
func NewTeamRepository(db *sqlx.DB) models.TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) GetAll() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Select(&teams, "SELECT id, name, attack, midfield, defence, overall FROM teams ORDER BY id")
	return teams, err
}

func (r *teamRepo) GetByID(id int) (*models.Team, error) {
	var team models.Team
	err := r.db.Get(&team, "SELECT id, name, attack, midfield, defence, overall FROM teams WHERE id = $1", id)
	return &team, err
}
