package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/domain"
	"github.com/searaaman/playledger/internal/services"
)

func CreatePayment(ctx *gin.Context) {
	var request domain.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var player domain.Player
	if err := config.DB.First(&player, request.PlayerID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "player not found",
		})
		return
	}

	session, ok := loadSession(ctx, request.SessionID)
	if !ok {
		return
	}

	playedInSession := false
	for _, bill := range services.CalculateSessionBills(session) {
		if bill.PlayerID == request.PlayerID {
			playedInSession = true
			break
		}
	}
	if !playedInSession {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "player is not part of this session",
		})
		return
	}

	payment := domain.Payment{
		PlayerID:  request.PlayerID,
		SessionID: request.SessionID,
		Amount:    request.Amount,
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to record payment",
		})
		return
	}
	payment.Player = &player
	ctx.JSON(http.StatusCreated, payment)
}

// ListPayments returns payments, newest first, optionally filtered by
// ?player_id= and/or ?session_id=.
func ListPayments(ctx *gin.Context) {
	query := config.DB.Preload("Player").Order("created_at DESC")
	if playerID := ctx.Query("player_id"); playerID != "" {
		query = query.Where("player_id = ?", playerID)
	}
	if sessionID := ctx.Query("session_id"); sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	payments := []domain.Payment{}
	if err := query.Find(&payments).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load payments",
		})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}

func DeletePayment(ctx *gin.Context) {
	paymentID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	result := config.DB.Delete(&domain.Payment{}, paymentID)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete payment",
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "payment not found",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}
