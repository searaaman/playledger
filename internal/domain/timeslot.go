package domain

import "time"

type TimeSlot struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	CourtsBooked int       `gorm:"column:courtsbooked" json:"courts_booked"`

	SessionID uint     `json:"session_id"`
	Players   []Player `gorm:"many2many:player_time_slots;" json:"players"`
}
