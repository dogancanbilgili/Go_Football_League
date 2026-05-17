package models

// Standing represents one row in the league table.
// It is computed from match results — it is never stored directly in the database.
// Team is embedded (struct composition): Standing automatically has all Team fields.
type Standing struct {
	Team
	Played         int     `json:"played"`
	Won            int     `json:"won"`
	Drawn          int     `json:"drawn"`
	Lost           int     `json:"lost"`
	GoalsFor       int     `json:"goals_for"`
	GoalsAgainst   int     `json:"goals_against"`
	GoalDifference int     `json:"goal_difference"`
	Points         int     `json:"points"`
}

// Prediction holds a team's estimated probability of winning the championship.
type Prediction struct {
	Team       Team    `json:"team"`
	Percentage float64 `json:"percentage"`
}
