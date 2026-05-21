package models

// Team represents a football club in the league.
// Attack, Midfield, Defence, and Overall come from EA FC ratings (sofifa.com).
type Team struct {
	ID       int    `db:"id"       json:"id"`
	Name     string `db:"name"     json:"name"`
	Attack   int    `db:"attack"   json:"attack"`
	Midfield int    `db:"midfield" json:"midfield"`
	Defence  int    `db:"defence"  json:"defence"`
	Overall  int    `db:"overall"  json:"overall"`
}
