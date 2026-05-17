-- =============================================================
-- TABLE: teams
-- Stores the 4 clubs participating in the league.
-- Strength (1-100) is used by the simulation to bias match results.
-- =============================================================
CREATE TABLE IF NOT EXISTS teams (
    id       SERIAL      PRIMARY KEY,
    name     VARCHAR(100) NOT NULL UNIQUE,
    strength INT          NOT NULL CHECK (strength BETWEEN 1 AND 100)
);

-- =============================================================
-- TABLE: matches
-- Stores all 12 fixtures for the season.
-- home_goals and away_goals are NULL until the match is simulated.
-- =============================================================
CREATE TABLE IF NOT EXISTS matches (
    id           SERIAL  PRIMARY KEY,
    week         INT     NOT NULL CHECK (week BETWEEN 1 AND 6),
    home_team_id INT     NOT NULL REFERENCES teams(id),
    away_team_id INT     NOT NULL REFERENCES teams(id),
    home_goals   INT     DEFAULT NULL CHECK (home_goals >= 0),
    away_goals   INT     DEFAULT NULL CHECK (away_goals >= 0),
    played       BOOLEAN NOT NULL DEFAULT false,

    -- A team cannot play itself
    CONSTRAINT different_teams CHECK (home_team_id <> away_team_id)
);

-- =============================================================
-- SEED: Insert the 4 teams
-- IDs will be assigned as 1, 2, 3, 4 by SERIAL
-- =============================================================
INSERT INTO teams (name, strength) VALUES
    ('Chelsea',          85),
    ('Arsenal',          82),
    ('Manchester City',  88),
    ('Liverpool',        84);

-- =============================================================
-- SEED: Insert all 12 fixtures (unplayed, no scores yet)
-- Week 1-3 = first half, Week 4-6 = return fixtures (home/away swapped)
-- =============================================================
INSERT INTO matches (week, home_team_id, away_team_id) VALUES
    -- Week 1
    (1, 1, 2),  -- Chelsea vs Arsenal
    (1, 3, 4),  -- Manchester City vs Liverpool

    -- Week 2
    (2, 1, 3),  -- Chelsea vs Manchester City
    (2, 2, 4),  -- Arsenal vs Liverpool

    -- Week 3
    (3, 1, 4),  -- Chelsea vs Liverpool
    (3, 2, 3),  -- Arsenal vs Manchester City

    -- Week 4 (return fixtures begin)
    (4, 2, 1),  -- Arsenal vs Chelsea
    (4, 4, 3),  -- Liverpool vs Manchester City

    -- Week 5
    (5, 3, 1),  -- Manchester City vs Chelsea
    (5, 4, 2),  -- Liverpool vs Arsenal

    -- Week 6
    (6, 4, 1),  -- Liverpool vs Chelsea
    (6, 3, 2);  -- Manchester City vs Arsenal
