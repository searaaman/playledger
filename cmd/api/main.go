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

	r.Run(":8080")



}