package models

// LeagueService defines what the business logic layer must be able to do.
// The handler layer depends on this interface, not on any concrete service type.
type LeagueService interface {
	GetTable() ([]Standing, error)
	GetWeekMatches(week int) ([]Match, error)
	SimulateNextWeek() ([]Match, error)
	SimulateAll() (map[int][]Match, error)
	GetPredictions() ([]Prediction, error)
	EditMatch(id, homeGoals, awayGoals int) error
	GetCurrentWeek() (int, error)
}
