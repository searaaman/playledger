package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/domain"
)

func CreatePlayer(ctx *gin.Context){
	var player domain.Player
	err:=ctx.ShouldBindJSON(&player)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return

	}
	err=config.DB.Create(&player).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated,player)
}

func GetPlayers(ctx *gin.Context){
	var players []domain.Player
	err:=config.DB.Find(&players).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK,players)

}
func AssignPlayerToTimeSlot(ctx *gin.Context){
	timeSlotID,err:=strconv.Atoi(ctx.Param("id"))
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid timeslot id",
		})
		return
	}

	var request domain.AssignPlayerRequest
	err=ctx.ShouldBindJSON(&request)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}

	var timeSlot domain.TimeSlot
	err=config.DB.First(&timeSlot,timeSlotID).Error
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"Error":"TimeSlot not found",
			})
		return
		
	}
	

	var player domain.Player
	err=config.DB.First(&player,request.PlayerID).Error
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"Error":"Player not found",
		})
		return

	}
	err=config.DB.
	Model(&timeSlot).
	Association("Players").
	Append(&player)

	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"Error":"Failed to assign player",

		})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"Message":"Player assigned succesfully",
	})

}