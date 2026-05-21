package database

import "github.com/jmoiron/sqlx"

// SeedFixtures generates and inserts the full round-robin schedule if the
// matches table is empty. It uses the circle method: one team is pinned,
// the rest rotate each round. After n-1 rounds every pair has met once
// (first half). The second half repeats with home and away swapped.
func SeedFixtures(db *sqlx.DB) error {
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM matches"); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Fetch team IDs in insertion order — order matters for fair scheduling
	var ids []int
	if err := db.Select(&ids, "SELECT id FROM teams ORDER BY id"); err != nil {
		return err
	}

	n := len(ids)
	totalRounds := n - 1

	type fixture struct{ week, home, away int }
	var fixtures []fixture

	// Work on a copy so we can rotate without touching the original slice
	working := make([]int, n)
	copy(working, ids)

	// First half: circle method
	for round := 0; round < totalRounds; round++ {
		week := round + 1
		for i := 0; i < n/2; i++ {
			fixtures = append(fixtures, fixture{week, working[i], working[n-1-i]})
		}
		// Rotate working[1:] right by one — pin working[0]
		last := working[n-1]
		copy(working[2:], working[1:n-1])
		working[1] = last
	}

	// Second half: same pairings, home and away swapped, weeks shifted
	firstHalf := len(fixtures)
	for i := 0; i < firstHalf; i++ {
		f := fixtures[i]
		fixtures = append(fixtures, fixture{f.week + totalRounds, f.away, f.home})
	}

	// Insert all fixtures in a single transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, f := range fixtures {
		if _, err := tx.Exec(
			"INSERT INTO matches (week, home_team_id, away_team_id) VALUES ($1, $2, $3)",
			f.week, f.home, f.away,
		); err != nil {
			tx.Rollback()// if there is an error, rollback the transaction.
			return err
		}
	}

	return tx.Commit()//only if the transaction is successful, commit the transaction.
}
