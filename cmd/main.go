package main

import (
	"log"

	"github.com/dogancanbilgili/Go_Football_League/internal/handler"
	"github.com/dogancanbilgili/Go_Football_League/internal/repository"
	"github.com/dogancanbilgili/Go_Football_League/internal/service"
	"github.com/dogancanbilgili/Go_Football_League/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Step 1: connect to PostgreSQL
	db, err := database.Connect()
	if err != nil {
		log.Fatal("could not connect to database:", err)
	}
	defer db.Close()

	// Step 2: seed fixtures if the matches table is empty
	if err := database.SeedFixtures(db); err != nil {
		log.Fatal("could not seed fixtures:", err)
	}

	// Step 3: wire the layers — each receives the interface of the layer below it
	teamRepo := repository.NewTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	svc := service.NewLeagueService(teamRepo, matchRepo)
	h := handler.NewLeagueHandler(svc)

	// Step 4: register routes
	router := gin.Default()

	// Allow browsers to call the API from any origin (required for the frontend)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	league := router.Group("/league")
	{
		league.GET("/table", h.GetTable)
		league.GET("/fixtures", h.GetFixtures)
		league.GET("/week/:week", h.GetWeekMatches)
		league.GET("/predictions", h.GetPredictions)
		league.POST("/next-week", h.SimulateNextWeek)
		league.POST("/play-all", h.SimulateAll)
		league.POST("/reset", h.Reset)
		league.PUT("/match/:id", h.EditMatch)
	}

	// Step 5: start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("could not start server:", err)
	}
}
