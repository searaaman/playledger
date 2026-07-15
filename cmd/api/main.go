package main

import (
	"github.com/searaaman/playledger/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/handlers"
	"github.com/searaaman/playledger/internal/middleware"
	"github.com/gin-contrib/cors"
)

func main(){
	config.ConnectDatabase()
	r:=gin.Default()

	r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:5173"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	AllowCredentials: true,
}))

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	r.GET("/health",handlers.HealthHandler) 
	auth.POST("/sessions",handlers.CreateSession)
	auth.POST("/sessions/:id/timeslots",handlers.CreateTimeSlot)
	auth.GET("/sessions/:id",handlers.GetSession)
	auth.POST("/players",handlers.CreatePlayer)
	auth.GET("/players",handlers.GetPlayers)
	auth.POST("/timeslots/:id/players",handlers.AssignPlayerToTimeSlot)
	auth.GET("/sessions/:id/billing",handlers.GetSessionBilling)
	auth.POST("/sessions/:id/billing",handlers.GetSessionBilling)
	auth.POST("/payments",handlers.CreatePayment)
	auth.GET("/players/:id/ledger",handlers.GetPlayerLedger)
	r.POST("/register",handlers.Register)
	r.POST("/login", handlers.Login)

	r.Run(":8080")



}