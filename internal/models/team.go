package models

// Team represents a football club in the league.
// Strength (1-100) biases match simulation: a stronger team scores more and concedes less.
type Team struct {
	ID       int    `db:"id"       json:"id"`
	Name     string `db:"name"     json:"name"`
	Strength int    `db:"strength" json:"strength"`
}
