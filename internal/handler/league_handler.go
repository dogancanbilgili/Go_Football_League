package handler

import (
	"net/http"
	"strconv"

	"github.com/dogancanbilgili/Go_Football_League/internal/models"
	"github.com/gin-gonic/gin"
)

type LeagueHandler struct {
	svc models.LeagueService
}

func NewLeagueHandler(svc models.LeagueService) *LeagueHandler {
	return &LeagueHandler{svc: svc}
}

// GET /league/table
func (h *LeagueHandler) GetTable(c *gin.Context) {
	standings, err := h.svc.GetTable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}

// GET /league/fixtures
func (h *LeagueHandler) GetFixtures(c *gin.Context) {
	matches, err := h.svc.GetFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GET /league/week/:week
func (h *LeagueHandler) GetWeekMatches(c *gin.Context) {
	week, err := strconv.Atoi(c.Param("week"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "week must be a number"})
		return
	}
	matches, err := h.svc.GetWeekMatches(week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// POST /league/next-week
func (h *LeagueHandler) SimulateNextWeek(c *gin.Context) {
	matches, err := h.svc.SimulateNextWeek()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// POST /league/play-all
func (h *LeagueHandler) SimulateAll(c *gin.Context) {
	results, err := h.svc.SimulateAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// PUT /league/match/:id
// Body: {"home_goals": 2, "away_goals": 1}
func (h *LeagueHandler) EditMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match id must be a number"})
		return
	}
	var body struct {
		HomeGoals int `json:"home_goals"`
		AwayGoals int `json:"away_goals"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.svc.EditMatch(id, body.HomeGoals, body.AwayGoals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match updated"})
}

// GET /league/predictions
func (h *LeagueHandler) GetPredictions(c *gin.Context) {
	predictions, err := h.svc.GetPredictions()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, predictions)
}

// POST /league/reset
func (h *LeagueHandler) Reset(c *gin.Context) {
	if err := h.svc.Reset(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "league reset to initial state"})
}
