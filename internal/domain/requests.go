package domain

import "time"

type CreatePlayerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
}

type CreateSessionRequest struct {
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
	CourtPrice float64   `json:"court_price" binding:"gte=0"`

	// When HourlySlots is true, the session is created with one slot per hour
	// between StartTime and EndTime, each booked with DefaultCourts courts.
	HourlySlots   bool `json:"hourly_slots"`
	DefaultCourts int  `json:"default_courts" binding:"gte=0"`
}

type CreateTimeSlotRequest struct {
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	CourtsBooked int       `json:"courts_booked" binding:"gte=0"`
}

type UpdateTimeSlotRequest struct {
	CourtsBooked int `json:"courts_booked" binding:"gte=0"`
}

type AssignPlayerRequest struct {
	PlayerID uint `json:"player_id" binding:"required"`
}

type CreatePaymentRequest struct {
	PlayerID  uint    `json:"player_id" binding:"required"`
	SessionID uint    `json:"session_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"gt=0"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
