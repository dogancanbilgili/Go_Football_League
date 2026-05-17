package models

// TeamRepository defines what the database layer must be able to do with teams.
// The service layer depends on this interface, not on any concrete database type.
type TeamRepository interface {
	GetAll() ([]Team, error)
	GetByID(id int) (*Team, error)
}

// MatchRepository defines what the database layer must be able to do with matches.
type MatchRepository interface {
	GetAll() ([]Match, error)
	GetByWeek(week int) ([]Match, error)
	GetCurrentWeek() (int, error)
	UpdateMatch(id, homeGoals, awayGoals int) error
	SaveSimulatedMatches(matches []Match) error
}
