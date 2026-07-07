package main

import (
	"github.com/searaaman/playledger/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/handlers"
)

func main(){
	config.ConnectDatabase()
	r:=gin.Default()

	r.GET("/health",handlers.HealthHandler) 
	r.POST("/sessions",handlers.CreateSession)
	r.POST("/sessions/:id/timeslots",handlers.CreateTimeSlot)
	r.GET("/sessions/:id",handlers.GetSession)
	r.POST("/players",handlers.CreatePlayer)
	r.GET("/players",handlers.GetPlayers)
	r.POST("/timeslots/:id/players",handlers.AssignPlayerToTimeSlot)
	r.GET("/sessions/:id/billing",handlers.GetSessionBilling)
	r.POST("/sessions/:id/billing",handlers.GetSessionBilling)
	r.POST("/payments",handlers.CreatePayment)

	r.Run(":8080")



}