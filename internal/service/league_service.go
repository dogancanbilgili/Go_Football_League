package service

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/dogancanbilgili/Go_Football_League/internal/models"
)

type leagueService struct {
	teamRepo  models.TeamRepository
	matchRepo models.MatchRepository
}

func NewLeagueService(teamRepo models.TeamRepository, matchRepo models.MatchRepository) models.LeagueService {
	return &leagueService{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
	}
}

// --- Simulation ---

// poissonGoals draws a random goal count using the Knuth algorithm.
// Multiply uniform randoms until their product falls below e^(-lambda).
// The number of draws follows a Poisson distribution with mean = lambda.
func poissonGoals(lambda float64) int {
	//lambda: the expected goal count
	L := math.Exp(-lambda) //e^(-lambda)
	k := 0 	//random draws count
	p := 1.0 //product of random numbers until p <= L, shrinks to 0 eventually
	for p > L {
		p *= rand.Float64()
		k++
	}
	return k - 1 //number of goals
}

// simulateMatch generates a scoreline using attack vs defence matchup.
// Home team gets a small advantage via higher lambda multiplier.
func simulateMatch(home, away models.Team) (homeGoals, awayGoals int) {
	homeStr := float64(home.Attack)*0.7 + float64(home.Midfield)*0.3
	awayStr := float64(away.Attack)*0.7 + float64(away.Midfield)*0.3
	ratioHome := homeStr / float64(away.Defence)
	ratioAway := awayStr / float64(home.Defence)
	lambdaHome := ratioHome * ratioHome * 1.4
	lambdaAway := ratioAway * ratioAway * 1.15
	return poissonGoals(lambdaHome), poissonGoals(lambdaAway)
}

// --- Table Calculation ---

func (s *leagueService) GetTable() ([]models.Standing, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	standingMap := make(map[int]*models.Standing)
	for _, t := range teams {
		t := t
		standingMap[t.ID] = &models.Standing{Team: t}
	}

	for _, m := range matches {
		if !m.Played {
			continue
		}
		home := standingMap[m.HomeTeamID]
		away := standingMap[m.AwayTeamID]
		hg, ag := *m.HomeGoals, *m.AwayGoals

		home.Played++
		away.Played++
		home.GoalsFor += hg
		home.GoalsAgainst += ag
		away.GoalsFor += ag
		away.GoalsAgainst += hg

		switch {
		case hg > ag:
			home.Won++
			home.Points += 3
			away.Lost++
		case ag > hg:
			away.Won++
			away.Points += 3
			home.Lost++
		default:
			home.Drawn++
			away.Drawn++
			home.Points++
			away.Points++
		}
	}

	standings := make([]models.Standing, 0, len(teams))
	for _, s := range standingMap {
		s.GoalDifference = s.GoalsFor - s.GoalsAgainst
		standings = append(standings, *s)
	}

	sort.Slice(standings, func(i, j int) bool {
		a, b := standings[i], standings[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDifference != b.GoalDifference {
			return a.GoalDifference > b.GoalDifference
		}
		if a.GoalsFor != b.GoalsFor {
			return a.GoalsFor > b.GoalsFor
		}
		return a.Name < b.Name
	})

	return standings, nil
}

// --- Week Simulation ---

func (s *leagueService) GetCurrentWeek() (int, error) {
	return s.matchRepo.GetCurrentWeek()
}

func (s *leagueService) GetWeekMatches(week int) ([]models.Match, error) {
	matches, err := s.matchRepo.GetByWeek(week)
	if err != nil {
		return nil, err
	}
	return s.attachTeams(matches)
}

func (s *leagueService) SimulateNextWeek() ([]models.Match, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	totalWeeks := 2 * (len(teams) - 1)

	week, err := s.matchRepo.GetCurrentWeek()
	if err != nil {
		return nil, err
	}
	if week > totalWeeks {
		return nil, fmt.Errorf("all %d weeks have already been played", totalWeeks)
	}

	matches, err := s.matchRepo.GetByWeek(week)
	if err != nil {
		return nil, err
	}

	teamMap := buildTeamMap(teams)
	for i, m := range matches {
		hg, ag := simulateMatch(teamMap[m.HomeTeamID], teamMap[m.AwayTeamID])
		matches[i].HomeGoals = &hg
		matches[i].AwayGoals = &ag
		matches[i].Played = true
	}

	if err := s.matchRepo.SaveSimulatedMatches(matches); err != nil {
		return nil, err
	}

	return s.attachTeams(matches)
}

func (s *leagueService) SimulateAll() (map[int][]models.Match, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	totalWeeks := 2 * (len(teams) - 1)

	week, err := s.matchRepo.GetCurrentWeek()
	if err != nil {
		return nil, err
	}

	results := make(map[int][]models.Match)
	for week <= totalWeeks {
		played, err := s.SimulateNextWeek()
		if err != nil {
			return nil, err
		}
		results[week] = played
		week++
	}
	return results, nil
}

// --- Championship Prediction ---

func (s *leagueService) GetPredictions() ([]models.Prediction, error) {
	week, err := s.matchRepo.GetCurrentWeek()
	if err != nil {
		return nil, err
	}
	// week is the next unplayed week, so completed weeks = week - 1
	if week < 5 {
		return nil, fmt.Errorf("predictions are available from week 4 onward (completed: week %d)", week-1)
	}

	standings, err := s.GetTable()
	if err != nil {
		return nil, err
	}

	allMatches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamMap := buildTeamMap(teams)

	scores := make(map[int]float64)
	for _, st := range standings {
		scores[st.ID] = float64(st.Points)
	}

	for _, m := range allMatches {
		if m.Played {
			continue
		}
		home := teamMap[m.HomeTeamID]
		away := teamMap[m.AwayTeamID]
		total := float64(home.Attack+away.Defence) + float64(away.Attack+home.Defence)

		homeWinProb := (float64(home.Attack+home.Midfield) / total) * 0.55
		awayWinProb := (float64(away.Attack+away.Midfield) / total) * 0.45
		drawProb := math.Max(0, 1.0-homeWinProb-awayWinProb)//cannot be negative

		scores[home.ID] += homeWinProb*3 + drawProb
		scores[away.ID] += awayWinProb*3 + drawProb
	}

	var totalScore float64
	for _, v := range scores {
		totalScore += v
	}

	predictions := make([]models.Prediction, 0, len(standings))
	for _, st := range standings {
		pct := 0.0
		if totalScore > 0 {
			pct = math.Round((scores[st.ID]/totalScore)*10000) / 100
		}
		predictions = append(predictions, models.Prediction{
			Team:       st.Team,
			Percentage: pct,
		})
	}

	return predictions, nil
}

// --- Fixtures ---

func (s *leagueService) GetFixtures() ([]models.Match, error) {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.attachTeams(matches)
}

// --- Reset ---

func (s *leagueService) Reset() error {
	return s.matchRepo.ResetMatches()
}

// --- Edit Match ---

func (s *leagueService) EditMatch(id, homeGoals, awayGoals int) error {
	if homeGoals < 0 || awayGoals < 0 {
		return fmt.Errorf("goals cannot be negative")
	}
	return s.matchRepo.UpdateMatch(id, homeGoals, awayGoals)
}

// --- Helpers ---

func buildTeamMap(teams []models.Team) map[int]models.Team {
	m := make(map[int]models.Team)
	for _, t := range teams {
		m[t.ID] = t
	}
	return m
}

func (s *leagueService) attachTeams(matches []models.Match) ([]models.Match, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamMap := buildTeamMap(teams)
	for i, m := range matches {
		home := teamMap[m.HomeTeamID]
		away := teamMap[m.AwayTeamID]
		matches[i].HomeTeam = &home
		matches[i].AwayTeam = &away
	}
	return matches, nil
}
