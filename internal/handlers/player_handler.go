package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searaaman/playledger/internal/config"
	"github.com/searaaman/playledger/internal/domain"
	"github.com/searaaman/playledger/internal/services"
)

func CreatePlayer(ctx *gin.Context) {
	var request domain.CreatePlayerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	player := domain.Player{
		Name:  request.Name,
		Phone: request.Phone,
	}
	if err := config.DB.Create(&player).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create player",
		})
		return
	}
	ctx.JSON(http.StatusCreated, player)
}

func GetPlayers(ctx *gin.Context) {
	players := []domain.Player{}
	if err := config.DB.Order("name").Find(&players).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load players",
		})
		return
	}
	ctx.JSON(http.StatusOK, players)
}

func AssignPlayerToTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var request domain.AssignPlayerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var timeSlot domain.TimeSlot
	if err := config.DB.First(&timeSlot, timeSlotID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "timeslot not found",
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

	// Appending an already-assigned player is a no-op, so this is safe to retry.
	if err := config.DB.Model(&timeSlot).Association("Players").Append(&player); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to assign player",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "player assigned successfully",
	})
}

func RemovePlayerFromTimeSlot(ctx *gin.Context) {
	timeSlotID, ok := parseID(ctx, "id")
	if !ok {
		return
	}
	playerID, ok := parseID(ctx, "playerId")
	if !ok {
		return
	}

	var timeSlot domain.TimeSlot
	if err := config.DB.First(&timeSlot, timeSlotID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "timeslot not found",
		})
		return
	}

	player := domain.Player{ID: playerID}
	if err := config.DB.Model(&timeSlot).Association("Players").Delete(&player); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove player",
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// loadLedgerData fetches everything the ledger calculation needs.
func loadLedgerData() ([]domain.Session, []domain.Payment, error) {
	var sessions []domain.Session
	if err := withSlots(config.DB).Find(&sessions).Error; err != nil {
		return nil, nil, err
	}
	var payments []domain.Payment
	if err := config.DB.Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, nil, err
	}
	return sessions, payments, nil
}

func GetLedger(ctx *gin.Context) {
	var players []domain.Player
	if err := config.DB.Find(&players).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load players",
		})
		return
	}
	sessions, payments, err := loadLedgerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load ledger",
		})
		return
	}
	ctx.JSON(http.StatusOK, services.BuildLedger(players, sessions, payments))
}

func GetPlayerLedger(ctx *gin.Context) {
	playerID, ok := parseID(ctx, "id")
	if !ok {
		return
	}

	var player domain.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "player not found",
		})
		return
	}

	sessions, payments, err := loadLedgerData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load ledger",
		})
		return
	}
	ctx.JSON(http.StatusOK, services.BuildPlayerLedger(player, sessions, payments))
}
