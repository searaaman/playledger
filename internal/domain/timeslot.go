package domain

import "time"

type TimeSlot struct{
	ID uint `gorm:"primaryKey"`
	StartTime time.Time
	EndTime time.Time
	Courtsbooked int
	players []*Player 
}

