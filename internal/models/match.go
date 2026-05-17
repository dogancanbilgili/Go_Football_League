package models

// Match represents a single fixture between two teams.
// HomeGoals and AwayGoals are pointers so they can be nil (match not yet played).
// A value of 0 means the team scored zero; nil means the match hasn't happened.
type Match struct {
	ID         int    `db:"id"           json:"id"`
	Week       int    `db:"week"         json:"week"`
	HomeTeamID int    `db:"home_team_id" json:"home_team_id"`
	AwayTeamID int    `db:"away_team_id" json:"away_team_id"`
	HomeGoals  *int   `db:"home_goals"   json:"home_goals"`
	AwayGoals  *int   `db:"away_goals"   json:"away_goals"`
	Played     bool   `db:"played"       json:"played"`
	HomeTeam   *Team  `db:"-"            json:"home_team,omitempty"`
	AwayTeam   *Team  `db:"-"            json:"away_team,omitempty"`
}
