package domain

import (
	"gorm.io/gorm"
)

type Payment struct{
	gorm.Model

	PlayerID uint
	SessionID uint
	Amount float64

	Player Player
	Session Session
}

type CreatePaymentRequest struct{
	PlayerID uint `json:"player_id"`
	SessionID uint `json:"session_id"`
	Amount float64 `json:"amount"`

}