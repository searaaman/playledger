package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/domain"

	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/services"
)


func CreatePayment(ctx *gin.Context){
	var request domain.CreatePaymentRequest
	err:=ctx.ShouldBindJSON(&request)

	if err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return
	}

	if request.Amount<0{

	}

	var player domain.Player 
	err=config.DB.First(&player ,request.PlayerID).Error
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"ERROR":"Player not found",
		})
		return
	}

	var session domain.Session
	err=config.DB.
	Preload("TimeSlots").
	Preload("TimeSlots.Players").
	First(&session,request.SessionID).Error
	if err!=nil{
		ctx.JSON(http.StatusNotFound,gin.H{
			"error":"session not found",
		})
		return
	}
	var playerBill *domain.PlayerBill
	bills:=services.CalculateSessionBills(session)
	for i:= range bills{
		if bills[i].PlayerID==request.PlayerID{
			playerBill=&bills[i]
			break
		}
	}
	if playerBill==nil{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"error":"Player is not part of this session",
		})
		return
	}

	payment:=domain.Payment{
	PlayerID:request.PlayerID,
	SessionID:request.SessionID,
	Amount:request.Amount,
	}
	err=config.DB.Create(&payment).Error
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated,payment)

}
