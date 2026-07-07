package handlers

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/domain"
)

func CreateSession(ctx *gin.Context){
	var session domain.Session
	err:=ctx.ShouldBindJSON(&session)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}
	err=config.DB.Create(&session).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return


	}
	ctx.JSON(http.StatusCreated,session)
	
}

func CreateTimeSlot(ctx *gin.Context){
	sessionID,err:=strconv.Atoi(ctx.Param("id"))
	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"Invalid Session ID",
		})
		return
	}

	var timeSlot domain.TimeSlot
	err=ctx.ShouldBindJSON(&timeSlot)
	if err !=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}
	timeSlot.SessionID=uint(sessionID)
	err=config.DB.Create(&timeSlot).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":"failed to create timeslot",
		})
		return
	}
	ctx.JSON(http.StatusCreated,timeSlot)
}
func GetSession(ctx *gin.Context){
	sessionID,err:=strconv.Atoi(ctx.Param("id"))
	if err!= nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"invalid session id ",
		})
		return
	}
	var session domain.Session
	err=config.DB.
		Preload("TimeSlots").
		Preload("TimeSlots.Players").
		First(&session,sessionID).Error
		if err!=nil{
			ctx.JSON(http.StatusNotFound,gin.H{
				"error":"Session not found",
			})
			return
		}
		ctx.JSON(http.StatusOK,session)
}

func GetSessionBilling(ctx *gin.Context){
	sessionID,err:=strconv.Atoi(ctx.Param("id"))
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"error":err.Error(),
		})
		return
	}
	

	var session domain.Session 
	err=config.DB.
		Preload("TimeSlots").
		Preload("TimeSlots.Players").
		First(&session,sessionID).Error
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"error":"Session not found",
		})
		return
	}
	bills:=CalculateSessionBills(session)
	ctx.JSON(http.StatusOK,bills)

}