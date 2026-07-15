package domain
import "time"

type Session struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	CourtPrice float64    `json:"court_price"`
	TimeSlots  []TimeSlot `json:"time_slots"`
}

