package domain

import "time"

type TimeSlot struct{
	ID uint `gorm:"primaryKey"`

	StartTime time.Time
	EndTime time.Time
	Courtsbooked int

	SessionID uint
	Players []Player `gorm:"many2many:player_time_slots;"`
}

