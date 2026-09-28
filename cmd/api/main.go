package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/handlers"
	"github.com/searaaman/playledger/internal/middleware"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}
	config.ConnectDatabase()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.HealthHandler)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/sessions", handlers.CreateSession)
	auth.POST("/sessions/:id/timeslots", handlers.CreateTimeSlot)
	auth.GET("/sessions/:id", handlers.GetSession)
	auth.POST("/players", handlers.CreatePlayer)
	auth.GET("/players", handlers.GetPlayers)
	auth.POST("/timeslots/:id/players", handlers.AssignPlayerToTimeSlot)
	auth.GET("/sessions/:id/billing", handlers.GetSessionBilling)
	auth.POST("/sessions/:id/billing", handlers.GetSessionBilling)
	auth.POST("/payments", handlers.CreatePayment)
	auth.GET("/players/:id/ledger", handlers.GetPlayerLedger)

	if err := r.Run(":" + config.Cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
